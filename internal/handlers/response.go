package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// PaginationMeta 分页元数据
type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// Success 成功响应
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Created 创建成功响应
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta 带元数据的成功响应
func SuccessWithMeta(c *gin.Context, message string, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// BadRequest 400错误响应
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

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: message,
	})
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Code:    http.StatusForbidden,
		Message: message,
	})
}

// NotFound 404错误响应
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

// Conflict 409错误响应
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Code:    http.StatusConflict,
		Message: message,
	})
}

// UnprocessableEntity 422错误响应
func UnprocessableEntity(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Errors:  errors,
	})
}

// TooManyRequests 429错误响应
func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, ErrorResponse{
		Code:    http.StatusTooManyRequests,
		Message: message,
	})
}

// InternalServerError 500错误响应
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

// ServiceUnavailable 503错误响应
func ServiceUnavailable(c *gin.Context, message string) {
	c.JSON(http.StatusServiceUnavailable, ErrorResponse{
		Code:    http.StatusServiceUnavailable,
		Message: message,
	})
}

// ValidationError 参数验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse 参数验证错误响应
func ValidationErrorResponse(c *gin.Context, errors []ValidationError) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: "参数验证失败",
		Errors:  errors,
	})
}

// PaginatedResponse 分页响应
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

// NoContent 204无内容响应
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Accepted 202已接受响应
func Accepted(c *gin.Context, message string) {
	c.JSON(http.StatusAccepted, Response{
		Code:    http.StatusAccepted,
		Message: message,
	})
}
