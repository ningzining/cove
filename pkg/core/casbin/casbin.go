package casbin

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

// Setup 设置 Casbin Enforcer
func Setup(db *gorm.DB) error {
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

// Enforcer 获取 Enforcer 实例
func Enforcer() *casbin.SyncedEnforcer {
	return enforcer
}
