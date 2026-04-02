package cmd

import (
	"github.com/ningzining/cove/app/sys/internal/config"
	"github.com/ningzining/cove/app/sys/internal/router"
	"github.com/ningzining/cove/pkg/core/conf"
	"github.com/ningzining/cove/pkg/rest"
	"github.com/ningzining/zlog"
	"github.com/spf13/cobra"
)

var configFile string

func Execute() error {
	rootCmd := &cobra.Command{
		// 指定命令的名字，该名字会出现在帮助信息中
		Use:   "cove-system",
		Short: "System Server",
		// 命令出错时，不打印帮助信息。设置为 true 可以确保命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// 指定调用 cmd.Execute() 时，执行的 Run 函数
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
	rootCmd.PersistentFlags().StringVarP(&configFile, "file", "f", "etc/auth.yaml", "Start server with provided configuration file")

	return rootCmd.Execute()
}

func run() {
	c := &config.Config{}
	conf.MustLoad(configFile, c)
	zlog.Init(&c.Log)

	server := rest.MustNewServer(&c.Config)
	router.MustRegister(server.Engine(), c)

	server.Start()
}
