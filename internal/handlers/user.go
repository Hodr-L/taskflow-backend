package handlers

import (
	"errors"
	"taskflow-backend/internal/models"
	"taskflow-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"taskflow-backend/internal/services"
)

type UserHandler struct {
	userService *services.UserService
	db          *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(db),
		db:          db,
	}
}

// getCurrentUserID 从上下文获取当前用户ID（UUID）
func (h *UserHandler) getCurrentUserID(c *gin.Context) (models.UUID, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return models.UUID{}, errors.New("未认证")
	}
	return models.ParseUUID(userIDStr.(string))
}

// parseUUIDParam 从路径参数解析UUID
func (h *UserHandler) parseUUIDParam(c *gin.Context, param string) (models.UUID, error) {
	idStr := c.Param(param)
	if idStr == "" {
		return models.UUID{}, errors.New("参数不能为空")
	}
	return models.ParseUUID(idStr)
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
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := h.getCurrentUserID(c)
	if err != nil {
		Unauthorized(c, err.Error())
		return
	}

	user, err := h.userService.GetUserByID(userID)
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
// @Param request body models.UpdateProfileRequest true "更新信息"
// @Success 200 {object} Response{data=models.UserResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := h.getCurrentUserID(c)
	if err != nil {
		Unauthorized(c, err.Error())
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数验证失败", err)
		return
	}

	user, err := h.userService.UpdateUser(userID, req)
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
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, err := h.getCurrentUserID(c)
	if err != nil {
		Unauthorized(c, err.Error())
		return
	}

	var req models.UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数验证失败", err)
		return
	}

	err = h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			NotFound(c, "用户不存在")
		case services.ErrInvalidCredentials:
			BadRequest(c, "原密码错误")
		case services.ErrInvalidInput:
			UnprocessableEntity(c, "新密码不符合要求", err)
		default:
			InternalServerError(c, "修改密码失败", err)
		}
		return
	}

	Success(c, "密码修改成功", nil)
}

// GetTeams 获取用户列表
func (h *UserHandler) GetListUsers(c *gin.Context) {
	var req models.GetUsersParams
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("查询参数验证失败", logger.ErrorField(err))
		BadRequest(c, "参数验证失败", err)
		return
	}

	listUser, total, totalPages, err := h.userService.ListUsers(req)

	if err != nil {
		InternalServerError(c, "查询失败", err)
		return
	}

	userData := h.userService.ToUserListResponse(listUser)
	response := models.UserListResponse{
		User: userData,
		Pagination: models.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	Success(c, "查询成功", response)

}

// GetUser GetTeam 获取用户详情
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := h.parseUUIDParam(c, "id")
	if err != nil {
		BadRequest(c, "无效的用户ID格式")
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			NotFound(c, "用户不存在")
		} else {
			InternalServerError(c, "获取用户信息失败", err)
		}
		return
	}

	Success(c, "获取成功", user.ToResponse())
}

// UpdateUser UpdateTeam 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := h.parseUUIDParam(c, "id")
	if err != nil {
		BadRequest(c, "无效的用户ID格式")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数验证失败", err)
		return
	}

	user, err := h.userService.UpdateUser(userID, req)
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

// DeleteUser DeleteTeam 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := h.parseUUIDParam(c, "id")
	if err != nil {
		BadRequest(c, "无效的用户ID格式")
		return
	}

	err = h.userService.DeleteUserByID(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			NotFound(c, "用户不存在")
		} else {
			InternalServerError(c, "删除失败", err)
		}
		return
	}

	Success(c, "删除成功", nil)
}

// ResetPassword 重置密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	userID, err := h.parseUUIDParam(c, "id")
	if err != nil {
		BadRequest(c, "无效的用户ID格式")
		return
	}

	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数验证失败", err)
		return
	}

	err = h.userService.ResetPassword(userID, req)
	if err != nil {
		InternalServerError(c, "重置密码错误", err)
		return
	}

	Success(c, "重置密码成功", nil)
}

// GetUsersStats 获取user 状态数量
func (h *UserHandler) GetUsersStats(c *gin.Context) {
	stats, err := h.userService.GetUserStatus()
	if err != nil {
		InternalServerError(c, "status 查询失败", err)
	}

	Success(c, "status查询成功", stats)
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数验证失败", err)
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		InternalServerError(c, "创建用户失败", err)
		return
	}
	Success(c, "创建用户成功", user.ToResponse())
}
