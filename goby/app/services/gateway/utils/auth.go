package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const RefreshTokenExpireDuration = 7 * 24 * time.Hour

var Secret = []byte("panic-bitdance")

func GenerateRefreshToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(RefreshTokenExpireDuration).Unix(), // 设置过期时间
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Secret) // 生成 JWT 格式的 refresh_token
}

// 判断 `refresh_token` 是否过期
func IsRefreshTokenExpired(refreshToken string) bool {
	// 解析 refresh_token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})

	if err != nil || !token.Valid {
		return true // token 解析失败，认为过期
	}

	// 从 token 中获取 Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return true
	}

	// 获取 `exp`（过期时间）
	exp, ok := claims["exp"].(float64) // JWT 的 `exp` 存储的是时间戳
	if !ok {
		return true
	}

	// 比较 `exp` 和当前时间
	return time.Now().Unix() > int64(exp) // 如果当前时间大于 `exp`，说明过期
}
