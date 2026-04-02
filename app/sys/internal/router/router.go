package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/app/sys/internal/config"
	"github.com/ningzining/cove/app/sys/internal/handler"
	"github.com/ningzining/cove/app/sys/internal/service"
	"github.com/ningzining/cove/app/sys/internal/svc"
	"github.com/ningzining/cove/pkg/model"
	"github.com/ningzining/cove/pkg/rest/response"
	"github.com/ningzining/cove/pkg/store"
	"github.com/ningzining/cove/pkg/xerr"
	"github.com/rs/zerolog/log"
)

func MustRegister(r *gin.Engine, c *config.Config) {
	// 注册健康检查路由
	r.GET("/healthz", func(ctx *gin.Context) {
		log.Info().Msg("healthz function called")
		response.Error(ctx, xerr.New(xerr.ErrCommon))
	})

	ctx := svc.NewContext(c)
	db := store.MustNew(&c.DB)

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal().Err(err).Msg("auto migrate user table failed")
	}

	authService := service.NewAuth(db, ctx)
	authHandler := handler.NewAuth(authService)

	// 注册路由
	api := r.Group("/api/v1")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)
	}
}
