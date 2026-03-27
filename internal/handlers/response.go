package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 缁熶竴鍝嶅簲缁撴瀯
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorResponse 閿欒鍝嶅簲缁撴瀯
type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// PaginationMeta 鍒嗛〉鍏冩暟鎹?type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// Success 鎴愬姛鍝嶅簲
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Created 鍒涘缓鎴愬姛鍝嶅簲
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta 甯﹀厓鏁版嵁鐨勬垚鍔熷搷搴?func SuccessWithMeta(c *gin.Context, message string, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// BadRequest 400閿欒鍝嶅簲
func BadRequest(c *gin.Context, message string, err ...error) {
	response := ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: message,
	}

	if len(err) > 0 && err[0] != nil {
		response.Errors = err[0].Error()
	}

	c.JSON(http.StatusBadRequest, response)
}

// Unauthorized 401閿欒鍝嶅簲
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: message,
	})
}

// Forbidden 403閿欒鍝嶅簲
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Code:    http.StatusForbidden,
		Message: message,
	})
}

// NotFound 404閿欒鍝嶅簲
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

// Conflict 409閿欒鍝嶅簲
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Code:    http.StatusConflict,
		Message: message,
	})
}

// UnprocessableEntity 422閿欒鍝嶅簲
func UnprocessableEntity(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Errors:  errors,
	})
}

// TooManyRequests 429閿欒鍝嶅簲
func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, ErrorResponse{
		Code:    http.StatusTooManyRequests,
		Message: message,
	})
}

// InternalServerError 500閿欒鍝嶅簲
func InternalServerError(c *gin.Context, message string, err ...error) {
	response := ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
	}

	if len(err) > 0 && err[0] != nil {
		response.Errors = err[0].Error()
	}

	c.JSON(http.StatusInternalServerError, response)
}

// ServiceUnavailable 503閿欒鍝嶅簲
func ServiceUnavailable(c *gin.Context, message string) {
	c.JSON(http.StatusServiceUnavailable, ErrorResponse{
		Code:    http.StatusServiceUnavailable,
		Message: message,
	})
}

// ValidationError 鍙傛暟楠岃瘉閿欒
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse 鍙傛暟楠岃瘉閿欒鍝嶅簲
func ValidationErrorResponse(c *gin.Context, errors []ValidationError) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: "鍙傛暟楠岃瘉澶辫触",
		Errors:  errors,
	})
}

// PaginatedResponse 鍒嗛〉鍝嶅簲
func PaginatedResponse(c *gin.Context, message string, data interface{}, page, limit int, total int64) {
	totalPage := int(total) / limit
	if int(total)%limit > 0 {
		totalPage++
	}

	meta := PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}

	SuccessWithMeta(c, message, data, meta)
}

// NoContent 204鏃犲唴瀹瑰搷搴?func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Accepted 202宸叉帴鍙楀搷搴?func Accepted(c *gin.Context, message string) {
	c.JSON(http.StatusAccepted, Response{
		Code:    http.StatusAccepted,
		Message: message,
	})
}
