package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/internal/pkg/response"
	"github.com/ningzining/cove/internal/pkg/xerr"
	"github.com/ningzining/cove/pkg/core/token"
	"github.com/rs/zerolog/log"
)

type contextKey string

const (
	ClaimsContextKey contextKey = "claims"
)

// AuthN JWT 认证中间件
func AuthN(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Error().Msg("no auth header")
			response.Error(c, xerr.New(xerr.ErrUnauthorized))
			c.Abort()
			return
		}

		// 从请求头中取出 token
		var tokenStr string
		fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)

		claims, err := token.Parse(tokenStr, key)
		if err != nil {
			log.Error().Err(err).Str("token", tokenStr).Msg("parse token failed")
			response.Error(c, xerr.New(xerr.ErrTokenInvalid))
			c.Abort()
			return
		}
		if claims.ExpiresAt.Before(time.Now()) {
			log.Error().Str("token", tokenStr).Msg("token expired")
			response.Error(c, xerr.New(xerr.ErrTokenExpired))
			c.Abort()
			return
		}

		c.Set(string(ClaimsContextKey), claims)

		ctx := context.WithValue(c.Request.Context(), ClaimsContextKey, claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetUserID 从 context 中获取当前用户的 UserID
func GetUserID(c *gin.Context) (string, bool) {
	claims, ok := c.Get(string(ClaimsContextKey))
	if !ok {
		return "", false
	}
	customClaims, ok := claims.(*token.CustomMapClaims)
	if !ok {
		return "", false
	}
	return customClaims.UserID, true
}

// MustGetUserID 从 context 中获取当前用户的 UserID，不存在则 panic
func MustGetUserID(c *gin.Context) string {
	userID, ok := GetUserID(c)
	if !ok {
		panic("user id not found in context")
	}
	return userID
}

