package middleware

import (
	"strings"
	"taskflow-backend/internal/services"

	"taskflow-backend/internal/handlers"
	"taskflow-backend/pkg/jwt"
	"taskflow-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWTAuth JWT璁よ瘉涓棿浠?func JWTAuth(jwtManager *jwt.JWTManager, tokenBlacklist *services.TokenBlacklistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 浠庤姹傚ご鑾峰彇token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handlers.Unauthorized(c, "缂哄皯璁よ瘉浠ょ墝")
			c.Abort()
			return
		}

		// 妫€鏌earer token鏍煎紡
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handlers.Unauthorized(c, "浠ょ墝鏍煎紡閿欒")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 1. 棣栧厛妫€鏌?token 鏄惁鍦ㄩ粦鍚嶅崟涓?		isBlacklisted, err := tokenBlacklist.IsBlacklisted(c, tokenString)
		if err != nil {
			logger.Error("妫€鏌ラ粦鍚嶅崟澶辫触", zap.Error(err))
			// Redis 閿欒鏃讹紝鍙互閫夋嫨缁х画楠岃瘉鎴栦笉楠岃瘉
			handlers.Unauthorized(c, "绯荤粺寮傚父")
			c.Abort()
			return
		}

		if isBlacklisted {
			logger.Warn("浣跨敤榛戝悕鍗曚腑鐨則oken", zap.String("token", tokenString[:20]+"..."))
			handlers.Unauthorized(c, "浠ょ墝宸插け鏁?)
			c.Abort()
			return
		}

		// 楠岃瘉token
		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			logger.Warn("JWT楠岃瘉澶辫触", logger.ErrorField(err), zap.String("path", c.Request.URL.Path))
			handlers.Unauthorized(c, "鏃犳晥鎴栬繃鏈熺殑浠ょ墝")
			c.Abort()
			return
		}

		// 灏嗙敤鎴蜂俊鎭瓨鍌ㄥ埌涓婁笅鏂?		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token", tokenString)

		logger.Debug("JWT璁よ瘉閫氳繃",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// OptionalJWTAuth 鍙€夌殑JWT璁よ瘉涓棿浠?func OptionalJWTAuth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
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
			// 浠ょ墝鏃犳晥锛屼絾涓嶉樆姝㈣姹?			logger.Debug("鍙€塉WT楠岃瘉澶辫触锛岀户缁鐞?, logger.ErrorField(err))
			c.Next()
			return
		}

		// 灏嗙敤鎴蜂俊鎭瓨鍌ㄥ埌涓婁笅鏂?		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token", tokenString)

		c.Next()
	}
}

// RequireRole 瑕佹眰鐗瑰畾瑙掕壊鐨勪腑闂翠欢
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			handlers.Unauthorized(c, "闇€瑕佽璇?)
			c.Abort()
			return
		}

		// 瑙掕壊鏉冮檺妫€鏌?		if !hasPermission(userRole.(string), requiredRole) {
			handlers.Forbidden(c, "鏉冮檺涓嶈冻")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin 瑕佹眰绠＄悊鍛樿鑹茬殑涓棿浠?func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireSuperAdmin 瑕佹眰瓒呯骇绠＄悊鍛樿鑹茬殑涓棿浠?func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole("super_admin")
}

// hasPermission 妫€鏌ョ敤鎴锋槸鍚︽湁鏉冮檺
func hasPermission(userRole, requiredRole string) bool {
	// 瑙掕壊鏉冮檺灞傜骇
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

// GetCurrentUserID 浠庝笂涓嬫枃鑾峰彇褰撳墠鐢ㄦ埛ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	if !ok {
		// 灏濊瘯杞崲涓篺loat64锛圝SON鏁板瓧鍙兘琚В鏋愪负float64锛?		if floatID, ok := userID.(float64); ok {
			return uint(floatID), true
		}
		return 0, false
	}

	return id, true
}

// GetCurrentUserRole 浠庝笂涓嬫枃鑾峰彇褰撳墠鐢ㄦ埛瑙掕壊
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

// IsAuthenticated 妫€鏌ョ敤鎴锋槸鍚﹀凡璁よ瘉
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// IsAdmin 妫€鏌ョ敤鎴锋槸鍚︽槸绠＄悊鍛?func IsAdmin(c *gin.Context) bool {
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

// IsSuperAdmin 妫€鏌ョ敤鎴锋槸鍚︽槸瓒呯骇绠＄悊鍛?func IsSuperAdmin(c *gin.Context) bool {
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
