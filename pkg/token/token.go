package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Key        string // 签名密钥
	ExpireTime int64  // 过期时间，单位秒
}

type CustomMapClaims struct {
	Provider string `json:"provider"` // 登录方式
	UserID   string `json:"user_id"`  // 用户ID
	Phone    string `json:"phone"`    // 手机号
	Nickname string `json:"nickname"` // 昵称
	jwt.RegisteredClaims
}

func Generate(claims jwt.Claims, signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Parse(tokenString string, signingKey string) (*CustomMapClaims, error) {
	claims := &CustomMapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
