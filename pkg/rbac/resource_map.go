package rbac

import "sync"

const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

const (
	ResourceUser   = "user"
	ResourceRole   = "role"
	ResourceTenant = "tenant"
)

type ResourceAction struct {
	Resource string
	Action   string
}

var (
	routeMap   = make(map[string]ResourceAction)
	routeMapMu sync.RWMutex
)

// RegisterRouteResource 注册路由与资源的映射关系
func RegisterRouteResource(method, path, resource, action string) {
	key := method + ":" + path
	routeMapMu.Lock()
	defer routeMapMu.Unlock()
	routeMap[key] = ResourceAction{
		Resource: resource,
		Action:   action,
	}
}

// GetResourceAction 获取路由对应的资源和操作
func GetResourceAction(method, path string) (ResourceAction, bool) {
	key := method + ":" + path
	routeMapMu.RLock()
	defer routeMapMu.RUnlock()
	ra, ok := routeMap[key]
	return ra, ok
}

// BatchRegisterRoutes 批量注册路由资源映射
func BatchRegisterRoutes(routes []RouteMapping) {
	for _, r := range routes {
		RegisterRouteResource(r.Method, r.Path, r.Resource, r.Action)
	}
}

type RouteMapping struct {
	Method   string
	Path     string
	Resource string
	Action   string
}

// GetDefaultRouteMappings 获取默认的路由资源映射
func GetDefaultRouteMappings() []RouteMapping {
	return []RouteMapping{
		// 用户管理
		{Method: "GET", Path: "/api/v1/users", Resource: ResourceUser, Action: ActionRead},
		{Method: "GET", Path: "/api/v1/users/:id", Resource: ResourceUser, Action: ActionRead},
		{Method: "POST", Path: "/api/v1/users", Resource: ResourceUser, Action: ActionCreate},
		{Method: "PUT", Path: "/api/v1/users/:id", Resource: ResourceUser, Action: ActionUpdate},
		{Method: "DELETE", Path: "/api/v1/users/:id", Resource: ResourceUser, Action: ActionDelete},

		// 角色管理
		{Method: "GET", Path: "/api/v1/roles", Resource: ResourceRole, Action: ActionRead},
		{Method: "GET", Path: "/api/v1/roles/:id", Resource: ResourceRole, Action: ActionRead},
		{Method: "POST", Path: "/api/v1/roles", Resource: ResourceRole, Action: ActionCreate},
		{Method: "PUT", Path: "/api/v1/roles/:id", Resource: ResourceRole, Action: ActionUpdate},
		{Method: "DELETE", Path: "/api/v1/roles/:id", Resource: ResourceRole, Action: ActionDelete},
	}
}
