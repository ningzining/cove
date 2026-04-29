package router

import (
	"errors"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/internal/system/config"
	"github.com/ningzining/cove/internal/system/handler"
	"github.com/ningzining/cove/internal/system/service"
	"github.com/ningzining/cove/internal/system/svc"
	"github.com/ningzining/cove/pkg/core/casbin"
	"github.com/ningzining/cove/pkg/model"
	"github.com/ningzining/cove/pkg/rest/middleware"
	"github.com/ningzining/cove/pkg/rest/response"
	"github.com/ningzining/cove/pkg/store"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func MustSetup(r *gin.Engine, cfg *config.Config) {
	db := store.MustNew(&cfg.DB)
	// 初始化casbin
	if err := casbin.Setup(db); err != nil {
		log.Fatal().Err(err).Msg("init casbin failed")
	}
	// 初始化数据库
	initDatabase(db)

	// 注册路由组
	g := r.Group("")
	// 注册系统路由
	setupSysRouter(g)
	// 注册业务路由
	setupBizRouter(g, cfg, db)
}

// setupSysRouter 注册系统路由
func setupSysRouter(g *gin.RouterGroup) {
	// 注册pprof路由
	pprof.Register(g)
	// 注册健康检查路由
	g.GET("/healthz", func(ctx *gin.Context) {
		log.Info().Msg("healthz function called")
		response.OK(ctx, nil)
	})
	// 仅在开发模式下 注册swagger路由
	if gin.Mode() != gin.ReleaseMode {
		g.GET("/swagger/system/*any", ginswagger.WrapHandler(swaggerfiles.NewHandler(), ginswagger.InstanceName("system")))
	}
	// 注册静态文件路由
	g.Static("/static", "./static")
}

// setupBizRouter 注册业务路由
func setupBizRouter(g *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	ctx := svc.NewContext(cfg)
	authService := service.NewAuth(db, ctx)
	authHandler := handler.NewAuth(authService)
	{
		v1 := g.Group("/api/v1")
		v1.POST("/login", authHandler.Login)
	}
	{
		roleService := service.NewRole(db, ctx)
		roleHandler := handler.NewRole(roleService)
		{
			v1 := g.Group("/api/v1/role")
			v1.Use(middleware.AuthN(&cfg.Jwt))
			v1.GET("/option", roleHandler.Option)
		}
		{
			v1 := g.Group("/api/v1/role")
			v1.Use(middleware.AuthN(&cfg.Jwt))
			v1.POST("", middleware.AuthZ(RoleResource, CreateAction), roleHandler.Create)
			v1.DELETE("", middleware.AuthZ(RoleResource, DeleteAction), roleHandler.Delete)
			v1.PUT("/:id", middleware.AuthZ(RoleResource, UpdateAction), roleHandler.Update)
			v1.GET("/:id", middleware.AuthZ(RoleResource, ReadAction), roleHandler.Get)
			v1.GET("", middleware.AuthZ(RoleResource, ReadAction), roleHandler.Page)
			v1.PUT("/:id/status", middleware.AuthZ(RoleResource, UpdateAction), roleHandler.UpdateStatus)
		}
	}
	{
		v1 := g.Group("/api/v1")
		v1.Use(middleware.AuthN(&cfg.Jwt))
	}
}

func initDatabase(db *gorm.DB) {
	migrate(db)
	initResources(db)
	initRoles(db)
	initUsers(db)
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Resource{},
		&model.UserRole{},
		&model.RoleResource{},
	); err != nil {
		log.Fatal().Err(err).Msg("auto migrate failed")
	}
}

func initResources(db *gorm.DB) {
	resources := []model.Resource{
		{Code: RoleResource, Name: "角色管理", Action: model.ReadAction},
		{Code: RoleResource, Name: "角色管理", Action: model.FullAction},
		{Code: UserResource, Name: "用户管理", Action: model.ReadAction},
		{Code: UserResource, Name: "用户管理", Action: model.FullAction},
	}
	for _, res := range resources {
		var count int64
		if err := db.Model(&model.Resource{}).Where("code = ? and action = ?", res.Code, res.Action).Count(&count).Error; err != nil {
			log.Fatal().Err(err).Msg("count resource failed")
		}
		if count == 0 {
			if err := db.Create(&res).Error; err != nil {
				log.Warn().Err(err).Str("code", res.Code).Msg("create resource failed")
			}
		}
	}
}

func initRoles(db *gorm.DB) {
	createRoleIfNotExists(db, model.AdminRoleCode, model.AdminRoleName, func() {
		enforcer := casbin.Enforcer()
		if enforcer == nil {
			log.Fatal().Msg("casbin enforcer not initialized")
			return
		}
		enforcer.AddPolicy(model.AdminRoleCode, "*", "*")
	})

	createRoleIfNotExists(db, model.NormalUserRoleCode, model.NormalUserRoleName, func() {
		enforcer := casbin.Enforcer()
		if enforcer == nil {
			log.Fatal().Msg("casbin enforcer not initialized")
			return
		}
		enforcer.AddPolicy(model.NormalUserRoleCode, "*", ReadAction)
	})
}

func createRoleIfNotExists(db *gorm.DB, code, name string, policy func()) {
	var role model.Role
	if err := db.Where("code = ?", code).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Str("code", code).Msg("create role")
			role = model.Role{Code: code, Name: name, Status: model.Enabled}
			if err := db.Create(&role).Error; err != nil {
				log.Fatal().Err(err).Str("code", code).Msg("create role failed")
			}
			policy()
		}
	}
}

func initUsers(db *gorm.DB) {
	createUserIfNotExists(db, model.AdminPhone, model.AdminNickname, model.AdminRoleCode)
	createUserIfNotExists(db, model.NormalUserPhone, model.NormalUserNickname, model.NormalUserRoleCode)
}

func createUserIfNotExists(db *gorm.DB, phone string, nickname, roleCode string) {
	var user model.User
	if err := db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Str("phone", phone).Msg("create user")
			user = model.User{
				Nickname: nickname,
				Phone:    phone,
				Password: model.DefaultPassword,
				Status:   model.Enabled,
			}
			if err := db.Create(&user).Error; err != nil {
				log.Fatal().Err(err).Str("phone", phone).Msg("create user failed")
			}
			// 创建数据库层面的用户-角色关联
			var role model.Role
			if err := db.Where("code = ?", roleCode).First(&role).Error; err != nil {
				log.Fatal().Err(err).Str("role_code", roleCode).Msg("find role failed")
			}
			userRole := model.UserRole{
				UserID: user.UserID,
				RoleID: role.RoleID,
			}
			if err := db.Create(&userRole).Error; err != nil {
				log.Fatal().Err(err).Str("phone", phone).Msg("create user role failed")
			}
			enforcer := casbin.Enforcer()
			if enforcer == nil {
				log.Fatal().Msg("casbin enforcer not initialized")
				return
			}
			enforcer.AddRoleForUser(user.UserID, roleCode)
		}
	}
}
