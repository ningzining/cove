package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/pkg/rest/middleware"
)

const (
	ActionCreate = "create"
	ActionDelete = "delete"
	ActionUpdate = "update"
	ActionRead   = "read"
)

const (
	ResourceUser = "user"
	ResourceRole = "role"
)

type ResourceAction struct {
	Resource string
	Action   string
}

func registerRouter(resource ResourceAction, group *gin.RouterGroup, method string, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return group.
		Use(middleware.AuthZ(resource.Resource, resource.Action)).
		Handle(method, relativePath, handlers...)
}
