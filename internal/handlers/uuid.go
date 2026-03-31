package handlers

import (
	"errors"
	"taskflow-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// parseUUIDParam 从路径参数解析UUID
func ParseUUIDParam(c *gin.Context, param string) (models.UUID, error) {
	idStr := c.Param(param)
	if idStr == "" {
		return models.UUID{}, errors.New("参数不能为空")
	}
	return models.ParseUUID(idStr)
}
