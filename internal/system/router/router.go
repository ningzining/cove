package router

import (
	"errors"
	"net/http"

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
	// 初始化数据库
	initDatabase(db)
	// 初始化casbin
	if err := casbin.Setup(db); err != nil {
		log.Fatal().Err(err).Msg("init casbin failed")
	}

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
		g.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.NewHandler()))
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
		v1 := g.Group("/api/v1")
		v1.Use(middleware.AuthN(&cfg.Jwt))
		registerRouter(ResourceAction{Resource: ResourceUser, Action: ActionRead}, v1, http.MethodGet, "/user", func(c *gin.Context) {
			response.OK(c, gin.H{"message": "用户列表"})
		})
		registerRouter(ResourceAction{Resource: ResourceUser, Action: ActionCreate}, v1, http.MethodPost, "/user", func(c *gin.Context) {
			response.OK(c, gin.H{"message": "创建用户"})
		})
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
		{Code: ResourceRole, Name: "角色管理", Actions: []string{ActionCreate, ActionRead, ActionUpdate, ActionDelete}},
		{Code: ResourceUser, Name: "用户管理", Actions: []string{ActionCreate, ActionRead, ActionUpdate, ActionDelete}},
	}
	for _, res := range resources {
		var count int64
		if err := db.Model(&model.Resource{}).Where("code = ?", res.Code).Count(&count).Error; err != nil {
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
	createRoleIfNotExists(db, model.RoleAdmin, "管理员")
	createRoleIfNotExists(db, model.RoleUser, "普通用户")
}

func createRoleIfNotExists(db *gorm.DB, code, name string) {
	var role model.Role
	if err := db.Where("code = ?", code).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Str("code", code).Msg("create role")
			role = model.Role{Code: code, Name: name, Status: model.Enabled}
			if err := db.Create(&role).Error; err != nil {
				log.Fatal().Err(err).Str("code", code).Msg("create role failed")
			}
		}
	}
}

func initUsers(db *gorm.DB) {
	createUserIfNotExists(db, model.AdminPhone, model.AdminNickname, model.RoleAdmin)
	createUserIfNotExists(db, model.NormalUserPhone, model.NormalUserNickname, model.RoleUser)
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
				UserID: user.ID,
				RoleID: role.ID,
			}
			if err := db.Create(&userRole).Error; err != nil {
				log.Fatal().Err(err).Str("phone", phone).Msg("create user role failed")
			}
		}
	}
}
