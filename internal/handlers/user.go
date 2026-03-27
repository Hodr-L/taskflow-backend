package handlers

import (
	"errors"
	"strconv"
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

// GetProfile 鑾峰彇褰撳墠鐢ㄦ埛淇℃伅
// @Summary 鑾峰彇褰撳墠鐢ㄦ埛淇℃伅
// @Description 鑾峰彇宸茬櫥褰曠敤鎴风殑璇︾粏淇℃伅
// @Tags 璁よ瘉
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Response{data=models.UserResponse}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /auth/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "鏈璇?)
		return
	}

	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		if err == services.ErrUserNotFound {
			NotFound(c, "鐢ㄦ埛涓嶅瓨鍦?)
		} else {
			InternalServerError(c, "鑾峰彇鐢ㄦ埛淇℃伅澶辫触", err)
		}
		return
	}

	Success(c, "鑾峰彇鎴愬姛", user.ToResponse())
}

// UpdateProfile 鏇存柊鐢ㄦ埛淇℃伅
// @Summary 鏇存柊鐢ㄦ埛淇℃伅
// @Description 鏇存柊褰撳墠鐢ㄦ埛鐨勪釜浜轰俊鎭?
// @Tags 璁よ瘉
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UpdateProfileRequest true "鏇存柊淇℃伅"
// @Success 200 {object} Response{data=models.UserResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "鏈璇?)
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "鍙傛暟楠岃瘉澶辫触", err)
		return
	}

	user, err := h.userService.UpdateUser(userID.(uint), req)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			NotFound(c, "鐢ㄦ埛涓嶅瓨鍦?)
		case services.ErrInvalidInput:
			BadRequest(c, "鏃犳晥鐨勮緭鍏?)
		default:
			InternalServerError(c, "鏇存柊鐢ㄦ埛淇℃伅澶辫触", err)
		}
		return
	}

	Success(c, "鏇存柊鎴愬姛", user.ToResponse())
}

// ChangePassword 淇敼瀵嗙爜
// @Summary 淇敼瀵嗙爜
// @Description 淇敼褰撳墠鐢ㄦ埛鐨勫瘑鐮?
// @Tags 璁よ瘉
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UserChangePasswordRequest true "瀵嗙爜淇℃伅"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "鏈璇?)
		return
	}

	var req models.UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "鍙傛暟楠岃瘉澶辫触", err)
		return
	}

	err := h.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			NotFound(c, "鐢ㄦ埛涓嶅瓨鍦?)
		case services.ErrInvalidCredentials:
			BadRequest(c, "鍘熷瘑鐮侀敊璇?)
		case services.ErrInvalidInput:
			UnprocessableEntity(c, "鏂板瘑鐮佷笉绗﹀悎瑕佹眰", err)
		default:
			InternalServerError(c, "淇敼瀵嗙爜澶辫触", err)
		}
		return
	}

	Success(c, "瀵嗙爜淇敼鎴愬姛", nil)
}

// GetTeams 鑾峰彇鐢ㄦ埛鍒楄〃
func (h *UserHandler) GetListUsers(c *gin.Context) {
	var req models.GetUsersParams
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("鏌ヨ鍙傛暟楠岃瘉澶辫触", logger.ErrorField(err))
		BadRequest(c, "鍙傛暟楠岃瘉澶辫触", err)
		return
	}

	listUser, total, totalPages, err := h.userService.ListUsers(req)

	if err != nil {
		InternalServerError(c, "鏌ヨ澶辫触", err)
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

	Success(c, "鏌ヨ鎴愬姛", response)

}

// GetUser GetTeam 鑾峰彇鐢ㄦ埛璇︽儏
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		InternalServerError(c, "鏃犳晥鐨勭敤鎴稩D")
		return
	}

	user, err := h.userService.GetUserByID(uint(userID))
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			NotFound(c, "鐢ㄦ埛涓嶅瓨鍦?)
		} else {
			InternalServerError(c, "鑾峰彇鐢ㄦ埛淇℃伅澶辫触", err)
		}
		return
	}

	Success(c, "鑾峰彇鎴愬姛", user.ToResponse())
}

// UpdateUser UpdateTeam 鏇存柊鐢ㄦ埛
func (h *UserHandler) UpdateUser(c *gin.Context) {

	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		InternalServerError(c, "鏃犳晥鐨勭敤鎴稩D")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "鍙傛暟楠岃瘉澶辫触", err)
		return
	}

	user, err := h.userService.UpdateUser(uint(userID), req)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			NotFound(c, "鐢ㄦ埛涓嶅瓨鍦?)
		case services.ErrInvalidInput:
			BadRequest(c, "鏃犳晥鐨勮緭鍏?)
		default:
			InternalServerError(c, "鏇存柊鐢ㄦ埛淇℃伅澶辫触", err)
		}
		return
	}

	Success(c, "鏇存柊鎴愬姛", user.ToResponse())
}

// DeleteUser DeleteTeam 鍒犻櫎鐢ㄦ埛
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// TODO: 瀹炵幇鍒犻櫎鐢ㄦ埛閫昏緫
	id := c.Param("id")

	userID, err2 := strconv.ParseUint(id, 10, 64)
	if err2 != nil {
		InternalServerError(c, "鏃犳晥鐨勭敤鎴稩D")
		return
	}

	err := h.userService.DeleteUserByID(uint(userID))
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			NotFound(c, "鐢ㄦ埛涓嶅瓨鍦?)
		} else {
			InternalServerError(c, "鍒犻櫎澶辫触", err)
		}
		return
	}

	Success(c, "鍒犻櫎鎴愬姛", nil)

}

// ResetPassword 閲嶇疆瀵嗙爜
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id := c.Param("id")

	userID, err2 := strconv.ParseUint(id, 10, 64)
	if err2 != nil {
		InternalServerError(c, "鏃犳晥鐨勭敤鎴稩D")
		return
	}

	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "鍙傛暟楠岃瘉澶辫触", err)
		return
	}

	err := h.userService.ResetPassword(uint(userID), req)
	if err != nil {
		InternalServerError(c, "閲嶇疆瀵嗙爜閿欒", err)
		return
	}

	Success(c, "閲嶇疆瀵嗙爜鎴愬姛", nil)

}

// GetUsersStats 鑾峰彇user 鐘舵€佹暟閲?
func (h *UserHandler) GetUsersStats(c *gin.Context) {
	stats, err := h.userService.GetUserStatus()
	if err != nil {
		InternalServerError(c, "status 鏌ヨ澶辫触", err)
	}

	Success(c, "status鏌ヨ鎴愬姛", stats)
}

// CreateUser 鍒涘缓鐢ㄦ埛
func (h *UserHandler) CreateUser(c *gin.Context) {

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "鍙傛暟楠岃瘉澶辫触", err)
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		InternalServerError(c, "鍒涘缓鐢ㄦ埛澶辫触", err)
		return
	}
	Success(c, "鍒涘缓鐢ㄦ埛鎴愬姛", user.ToResponse())
}
