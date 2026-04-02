package svc

import "github.com/ningzining/cove/app/sys/internal/config"

type Context struct {
	Config *config.Config
}

func NewContext(config *config.Config) *Context {
	return &Context{
		Config: config,
	}
}
