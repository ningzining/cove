package rbac

import (
	"sync"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const modelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && (p.obj == "*" || r.obj == p.obj) && (p.act == "*" || r.act == p.act)
`

var (
	enforcer *casbin.SyncedEnforcer
	once     sync.Once
)

// Init 初始化 Casbin Enforcer
func Init(db *gorm.DB) error {
	var err error
	once.Do(func() {
		var m model.Model
		m, err = model.NewModelFromString(modelText)
		if err != nil {
			err = errors.Wrap(err, "failed to load model from text")
			return
		}

		var adapter *gormadapter.Adapter
		adapter, err = gormadapter.NewAdapterByDBUseTableName(db, "sys", "casbin_rule")
		if err != nil {
			err = errors.Wrap(err, "failed to create gorm adapter")
			return
		}

		enforcer, err = casbin.NewSyncedEnforcer(m, adapter)
		if err != nil {
			err = errors.Wrap(err, "failed to create enforcer")
			return
		}

		if err = enforcer.LoadPolicy(); err != nil {
			err = errors.Wrap(err, "failed to load policy")
			return
		}
	})
	return err
}

// GetEnforcer 获取 Enforcer 实例
func GetEnforcer() *casbin.SyncedEnforcer {
	return enforcer
}

// CheckPermission 检查权限
func CheckPermission(userID, resource, action string) (bool, error) {
	return enforcer.Enforce(userID, resource, action)
}

// AddRoleForUser 为用户添加角色
func AddRoleForUser(userID, roleCode string) (bool, error) {
	return enforcer.AddRoleForUser(userID, roleCode)
}

// DeleteRoleForUser 删除用户角色
func DeleteRoleForUser(userID, roleCode string) (bool, error) {
	return enforcer.DeleteRoleForUser(userID, roleCode)
}

// GetRolesForUser 获取用户的角色
func GetRolesForUser(userID string) ([]string, error) {
	return enforcer.GetRolesForUser(userID)
}

// AddPolicy 添加策略
func AddPolicy(roleCode, resource, action string) (bool, error) {
	return enforcer.AddPolicy(roleCode, resource, action)
}

// RemovePolicy 删除策略
func RemovePolicy(roleCode, resource, action string) (bool, error) {
	return enforcer.RemovePolicy(roleCode, resource, action)
}

// LoadPolicy 重新加载策略
func LoadPolicy() error {
	return enforcer.LoadPolicy()
}
