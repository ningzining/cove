package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/internal/pkg/response"
	"github.com/ningzining/cove/internal/pkg/xerr"
	"github.com/ningzining/cove/pkg/core/casbin"
	"github.com/ningzining/cove/pkg/core/token"
	"github.com/rs/zerolog/log"
)

// AuthZ Casbin 权限检查中间件
func AuthZ(resource string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims, ok := getClaimsFromContext(c)
		if !ok {
			log.Error().Msg("claims not found in context")
			response.Error(c, xerr.New(xerr.ErrUnauthorized))
			c.Abort()
			return
		}
		allowed, err := casbin.Enforcer().Enforce(userClaims.UserID, resource, action)
		if err != nil {
			log.Error().Err(err).Msg("check permission failed")
			response.Error(c, xerr.New(xerr.ErrForbidden))
			c.Abort()
			return
		}

		if !allowed {
			log.Warn().
				Str("path", c.FullPath()).
				Str("user_id", userClaims.UserID).
				Str("resource", resource).
				Str("action", action).
				Msg("permission denied")
			response.Error(c, xerr.New(xerr.ErrForbidden))
			c.Abort()
			return
		}

		c.Next()
	}
}

func getClaimsFromContext(c *gin.Context) (*token.CustomMapClaims, bool) {
	claims, ok1 := c.Get(string(ClaimsContextKey))
	if !ok1 {
		return nil, false
	}
	cl, ok := claims.(*token.CustomMapClaims)
	return cl, ok
}
