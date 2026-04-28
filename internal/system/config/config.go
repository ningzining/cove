package config

import (
	"github.com/ningzining/cove/pkg/core/zlog"
	"github.com/ningzining/cove/pkg/rest"
	"github.com/ningzining/cove/pkg/store"
	"github.com/ningzining/cove/pkg/token"
)

type Config struct {
	rest.Config `mapstructure:",squash"` // 应用配置

	Jwt token.Config `mapstructure:"jwt"` // JWT配置
	Log zlog.Config  `mapstructure:"log"` // 日志配置
	DB  store.Config `mapstructure:"db"`  // 数据库配置
}
