package token

import (
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Key        string // 签名密钥
	ExpireTime int64  // 过期时间，单位秒
}

type CustomMapClaims struct {
	Provider string `json:"provider"` // 登录方式
	UserID   int64  `json:"user_id"`  // 用户ID
	Phone    string `json:"phone"`    // 手机号
	Nickname string `json:"nickname"` // 昵称
	jwt.Claims
}

func Generate(claims jwt.Claims, signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
