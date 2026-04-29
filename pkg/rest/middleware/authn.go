package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/pkg/rest/response"
	"github.com/ningzining/cove/pkg/token"
	"github.com/ningzining/cove/pkg/xerr"
	"github.com/rs/zerolog/log"
)

type contextKey string

const (
	ClaimsContextKey contextKey = "claims"
)

// AuthN JWT 认证中间件
func AuthN(cfg *token.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg == nil {
			log.Error().Msg("token config not set")
			response.Error(c, xerr.New(xerr.ErrUnauthorized))
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Error().Msg("no auth header")
			response.Error(c, xerr.New(xerr.ErrUnauthorized))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			log.Error().Msg("invalid auth header")
			response.Error(c, xerr.New(xerr.ErrUnauthorized))
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := token.Parse(tokenStr, cfg.Key)
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

func GetClaimsFromContext(c *gin.Context) (*token.CustomMapClaims, bool) {
	claims, ok1 := c.Get(string(ClaimsContextKey))
	if !ok1 {
		return nil, false
	}
	cl, ok := claims.(*token.CustomMapClaims)
	return cl, ok
}
