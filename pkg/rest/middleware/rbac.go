package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/pkg/rbac"
	"github.com/ningzining/cove/pkg/rest/response"
	"github.com/ningzining/cove/pkg/xerr"
	"github.com/rs/zerolog/log"
)

// RBAC Casbin 权限检查中间件
func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims, ok := GetClaimsFromContext(c)
		if !ok {
			response.Error(c, xerr.New(xerr.ErrUnauthorized))
			c.Abort()
			return
		}
		userID := userClaims.UserID

		method := c.Request.Method
		path := c.FullPath()

		ra, ok := rbac.GetResourceAction(method, path)
		if !ok {
			c.Next()
			return
		}

		enforcer := rbac.GetEnforcer()
		if enforcer == nil {
			log.Error().Msg("casbin enforcer not initialized")
			response.Error(c, xerr.New(xerr.ErrForbidden))
			c.Abort()
			return
		}

		allowed, err := enforcer.Enforce(userID, ra.Resource, ra.Action)
		if err != nil {
			log.Error().Err(err).Msg("check permission failed")
			response.Error(c, xerr.New(xerr.ErrForbidden))
			c.Abort()
			return
		}

		if !allowed {
			log.Warn().
				Str("path", path).
				Str("user_id", userID).
				Str("resource", ra.Resource).
				Str("action", ra.Action).
				Msg("permission denied")
			response.Error(c, xerr.New(xerr.ErrForbidden))
			c.Abort()
			return
		}

		c.Next()
	}
}
