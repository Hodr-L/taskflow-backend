package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"taskflow-backend/internal/models"
	"taskflow-backend/internal/services"
	"taskflow-backend/pkg/jwt"
	"taskflow-backend/pkg/logger"
)

type AuthHandler struct {
	authService *services.AuthService
	jwtManager  *jwt.JWTManager
	db          *gorm.DB
}

func NewAuthHandler(db *gorm.DB, jwtManager *jwt.JWTManager) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(db),
		jwtManager:  jwtManager,
		db:          db,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.UserCreateRequest true "注册信息"
// @Success 201 {object} Response{data=models.AuthResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("注册参数验证失败", logger.ErrorField(err))
		BadRequest(c, "参数验证失败", err)
		return
	}

	// 创建用户
	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		switch err {
		case services.ErrUserExists:
			Conflict(c, "用户已存在")
		case services.ErrInvalidInput:
			BadRequest(c, "无效的输入")
		default:
			InternalServerError(c, "注册失败", err)
		}
		return
	}

	// 生成JWT令牌
	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user)
	if err != nil {
		InternalServerError(c, "生成令牌失败", err)
		return
	}

	// 更新最后登录时间
	user.UpdateLastLogin()
	h.db.Save(user)

	response := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(24 * time.Hour.Seconds()), // 24小时
		User:         user.ToResponse(),
	}

	Created(c, "注册成功", response)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.UserLoginRequest true "登录信息"
// @Success 200 {object} Response{data=models.AuthResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("登录参数验证失败", logger.ErrorField(err))
		BadRequest(c, "参数验证失败", err)
		return
	}

	// 验证用户
	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			NotFound(c, "用户不存在")
		case services.ErrInvalidCredentials:
			Unauthorized(c, "用户名或密码错误")
		case services.ErrUserInactive:
			Unauthorized(c, "账户已被禁用")
		default:
			InternalServerError(c, "登录失败", err)
		}
		return
	}

	// 生成JWT令牌
	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user)
	if err != nil {
		InternalServerError(c, "生成令牌失败", err)
		return
	}

	// 更新最后登录时间
	user.UpdateLastLogin()
	h.db.Save(user)

	response := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(24 * time.Hour.Seconds()),
		User:         user.ToResponse(),
	}

	Success(c, "登录成功", response)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出（客户端应删除令牌）
// @Tags 认证
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Response
// @Failure 401 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 在实际应用中，可以将令牌加入黑名单
	// 这里只是返回成功，由客户端删除令牌
	Success(c, "登出成功", nil)
}

// RefreshToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新令牌"
// @Success 200 {object} Response{data=RefreshTokenResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数验证失败", err)
		return
	}

	accessToken, claims, err := h.jwtManager.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		switch err {
		case jwt.ErrExpiredToken:
			Unauthorized(c, "刷新令牌已过期")
		case jwt.ErrInvalidToken:
			Unauthorized(c, "无效的刷新令牌")
		default:
			InternalServerError(c, "刷新令牌失败", err)
		}
		return
	}

	response := RefreshTokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(24 * time.Hour.Seconds()),
		UserID:      claims.UserID,
		Username:    claims.Username,
	}

	Success(c, "令牌刷新成功", response)
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取已登录用户的详细信息
// @Tags 认证
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Response{data=models.UserResponse}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未认证")
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		if err == services.ErrUserNotFound {
			NotFound(c, "用户不存在")
		} else {
			InternalServerError(c, "获取用户信息失败", err)
		}
		return
	}

	Success(c, "获取成功", user.ToResponse())
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前用户的个人信息
// @Tags 认证
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UserUpdateRequest true "更新信息"
// @Success 200 {object} Response{data=models.UserResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未认证")
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数验证失败", err)
		return
	}

	user, err := h.authService.UpdateUser(userID.(uint), req)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			NotFound(c, "用户不存在")
		case services.ErrInvalidInput:
			BadRequest(c, "无效的输入")
		default:
			InternalServerError(c, "更新用户信息失败", err)
		}
		return
	}

	Success(c, "更新成功", user.ToResponse())
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的密码
// @Tags 认证
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UserChangePasswordRequest true "密码信息"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "未认证")
		return
	}

	var req models.UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数验证失败", err)
		return
	}

	err := h.authService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			NotFound(c, "用户不存在")
		case services.ErrInvalidCredentials:
			Unauthorized(c, "原密码错误")
		case services.ErrInvalidInput:
			BadRequest(c, "新密码不符合要求")
		default:
			InternalServerError(c, "修改密码失败", err)
		}
		return
	}

	Success(c, "密码修改成功", nil)
}

// 请求和响应结构体
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
}
