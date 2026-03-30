package jwt

import (
	"errors"
	"fmt"
	"time"

	"taskflow-backend/internal/config"
	"taskflow-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("无效的token")
	ErrExpiredToken = errors.New("token已过期")
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTManager(cfg config.JWTConfig) *JWTManager {
	return &JWTManager{
		secretKey:            cfg.Secret,
		accessTokenDuration:  time.Hour * time.Duration(cfg.AccessTokenExpireHours),
		refreshTokenDuration: time.Hour * 24 * time.Duration(cfg.RefreshTokenExpireDays),
	}
}

// GenerateAccessToken 生成访问令牌
func (m *JWTManager) GenerateAccessToken(user *models.User) (string, error) {
	claims := Claims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "taskflow-backend",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// GenerateRefreshToken 生成刷新令牌
func (m *JWTManager) GenerateRefreshToken(user *models.User) (string, error) {
	claims := Claims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "taskflow-backend",
			Subject:   fmt.Sprintf("refresh-%s", user.ID.String()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// VerifyToken 验证令牌
func (m *JWTManager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshAccessToken 刷新访问令牌
func (m *JWTManager) RefreshAccessToken(refreshToken string) (string, *Claims, error) {
	claims, err := m.VerifyToken(refreshToken)
	if err != nil {
		return "", nil, err
	}

	// 检查是否是刷新令牌
	if claims.Subject != fmt.Sprintf("refresh-%s", claims.UserID) {
		return "", nil, ErrInvalidToken
	}

	// 生成新的访问令牌
	newClaims := Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "taskflow-backend",
			Subject:   claims.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	accessToken, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", nil, err
	}

	return accessToken, claims, nil
}

// GenerateTokens 生成访问令牌和刷新令牌
func (m *JWTManager) GenerateTokens(user *models.User) (accessToken, refreshToken string, err error) {
	accessToken, err = m.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = m.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析令牌（不验证过期）
func (m *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GetTokenDuration 获取令牌剩余时间
func (m *JWTManager) GetTokenDuration(tokenString string) (time.Duration, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	if claims.ExpiresAt == nil {
		return 0, errors.New("token没有过期时间")
	}

	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining < 0 {
		return 0, ErrExpiredToken
	}

	return remaining, nil
}
