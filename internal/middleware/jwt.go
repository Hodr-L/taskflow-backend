package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"taskflow-backend/internal/handlers"
	"taskflow-backend/pkg/jwt"
	"taskflow-backend/pkg/logger"
)

// JWTAuth JWT认证中间件
func JWTAuth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handlers.Unauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handlers.Unauthorized(c, "令牌格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证token
		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			logger.Warn("JWT验证失败", logger.ErrorField(err), zap.String("path", c.Request.URL.Path))
			handlers.Unauthorized(c, "无效或过期的令牌")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token", tokenString)

		logger.Debug("JWT认证通过",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// OptionalJWTAuth 可选的JWT认证中间件
func OptionalJWTAuth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			// 令牌无效，但不阻止请求
			logger.Debug("可选JWT验证失败，继续处理", logger.ErrorField(err))
			c.Next()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token", tokenString)

		c.Next()
	}
}

// RequireRole 要求特定角色的中间件
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			handlers.Unauthorized(c, "需要认证")
			c.Abort()
			return
		}

		// 角色权限检查
		if !hasPermission(userRole.(string), requiredRole) {
			handlers.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin 要求管理员角色的中间件
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireSuperAdmin 要求超级管理员角色的中间件
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole("super_admin")
}

// hasPermission 检查用户是否有权限
func hasPermission(userRole, requiredRole string) bool {
	// 角色权限层级
	roleHierarchy := map[string]int{
		"user":        1,
		"admin":       2,
		"super_admin": 3,
	}

	userLevel, userOk := roleHierarchy[userRole]
	requiredLevel, requiredOk := roleHierarchy[requiredRole]

	if !userOk || !requiredOk {
		return false
	}

	return userLevel >= requiredLevel
}

// GetCurrentUserID 从上下文获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	if !ok {
		// 尝试转换为float64（JSON数字可能被解析为float64）
		if floatID, ok := userID.(float64); ok {
			return uint(floatID), true
		}
		return 0, false
	}

	return id, true
}

// GetCurrentUserRole 从上下文获取当前用户角色
func GetCurrentUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}

	roleStr, ok := role.(string)
	if !ok {
		return "", false
	}

	return roleStr, true
}

// IsAuthenticated 检查用户是否已认证
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// IsAdmin 检查用户是否是管理员
func IsAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}

	roleStr, ok := role.(string)
	if !ok {
		return false
	}

	return roleStr == "admin" || roleStr == "super_admin"
}

// IsSuperAdmin 检查用户是否是超级管理员
func IsSuperAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}

	roleStr, ok := role.(string)
	if !ok {
		return false
	}

	return roleStr == "super_admin"
}
