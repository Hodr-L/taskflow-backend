package handlers

import (
	"time"

	"taskflow-backend/internal/models"
	"taskflow-backend/internal/services"
	"taskflow-backend/pkg/jwt"
	"taskflow-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// Register 鐢ㄦ埛娉ㄥ唽
// @Summary 鐢ㄦ埛娉ㄥ唽
// @Description 鍒涘缓鏂扮敤鎴疯处鎴?// @Tags 璁よ瘉
// @Accept json
// @Produce json
// @Param request body models.UserCreateRequest true "娉ㄥ唽淇℃伅"
// @Success 201 {object} Response{data=models.AuthResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("娉ㄥ唽鍙傛暟楠岃瘉澶辫触", logger.ErrorField(err))
		BadRequest(c, "鍙傛暟楠岃瘉澶辫触", err)
		return
	}

	// 鍒涘缓鐢ㄦ埛
	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		switch err {
		case services.ErrUserExists:
			Conflict(c, "鐢ㄦ埛宸插瓨鍦?)
		case services.ErrInvalidInput:
			BadRequest(c, "鏃犳晥鐨勮緭鍏?)
		default:
			InternalServerError(c, "娉ㄥ唽澶辫触", err)
		}
		return
	}

	// 鐢熸垚JWT浠ょ墝
	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user)
	if err != nil {
		InternalServerError(c, "鐢熸垚浠ょ墝澶辫触", err)
		return
	}

	// 鏇存柊鏈€鍚庣櫥褰曟椂闂?	user.UpdateLastLogin()
	h.db.Save(user)

	response := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(24 * time.Hour.Seconds()), // 24灏忔椂
		User:         user.ToResponse(),
	}

	Created(c, "鎴愬姛", response)
}

// Login 鐢ㄦ埛鐧诲綍
// @Summary 鐢ㄦ埛鐧诲綍
// @Description 鐢ㄦ埛鐧诲綍鑾峰彇璁块棶浠ょ墝
// @Tags 璁よ瘉
// @Accept json
// @Produce json
// @Param request body models.UserLoginRequest true "鐧诲綍淇℃伅"
// @Success 200 {object} Response{data=models.AuthResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("鐧诲綍鍙傛暟楠岃瘉澶辫触", logger.ErrorField(err))
		BadRequest(c, "鍙傛暟楠岃瘉澶辫触", err)
		return
	}

	// 楠岃瘉鐢ㄦ埛
	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			NotFound(c, "鐢ㄦ埛涓嶅瓨鍦?)
		case services.ErrInvalidCredentials:
			Unauthorized(c, "鐢ㄦ埛鍚嶆垨瀵嗙爜閿欒")
		case services.ErrUserInactive:
			Unauthorized(c, "璐︽埛宸茶绂佺敤")
		default:
			InternalServerError(c, "鐧诲綍澶辫触", err)
		}
		return
	}

	// 鐢熸垚JWT浠ょ墝
	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user)
	if err != nil {
		InternalServerError(c, "鐢熸垚浠ょ墝澶辫触", err)
		return
	}

	// 鏇存柊鏈€鍚庣櫥褰曟椂闂?	user.UpdateLastLogin()
	h.db.Save(user)

	response := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(24 * time.Hour.Seconds()),
		User:         user.ToResponse(),
	}

	Success(c, "鐧诲綍鎴愬姛", response)
}

// RefreshToken 鍒锋柊璁块棶浠ょ墝
// @Summary 鍒锋柊璁块棶浠ょ墝
// @Description 浣跨敤鍒锋柊浠ょ墝鑾峰彇鏂扮殑璁块棶浠ょ墝
// @Tags 璁よ瘉
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "鍒锋柊浠ょ墝"
// @Success 200 {object} Response{data=RefreshTokenResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "鍙傛暟楠岃瘉澶辫触", err)
		return
	}

	accessToken, claims, err := h.jwtManager.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		switch err {
		case jwt.ErrExpiredToken:
			Unauthorized(c, "鍒锋柊浠ょ墝宸茶繃鏈?)
		case jwt.ErrInvalidToken:
			Unauthorized(c, "鏃犳晥鐨勫埛鏂颁护鐗?)
		default:
			InternalServerError(c, "鍒锋柊浠ょ墝澶辫触", err)
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

	Success(c, "浠ょ墝鍒锋柊鎴愬姛", response)
}

// 璇锋眰鍜屽搷搴旂粨鏋勪綋
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
