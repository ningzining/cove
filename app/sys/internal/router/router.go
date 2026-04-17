package router

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/app/sys/internal/config"
	"github.com/ningzining/cove/app/sys/internal/handler"
	"github.com/ningzining/cove/app/sys/internal/service"
	"github.com/ningzining/cove/app/sys/internal/svc"
	"github.com/ningzining/cove/pkg/model"
	"github.com/ningzining/cove/pkg/rbac"
	"github.com/ningzining/cove/pkg/rest/middleware"
	"github.com/ningzining/cove/pkg/rest/response"
	"github.com/ningzining/cove/pkg/store"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func MustRegister(r *gin.Engine, c *config.Config) {
	r.GET("/healthz", func(ctx *gin.Context) {
		log.Info().Msg("healthz function called")
		response.OK(ctx, nil)
	})

	ctx := svc.NewContext(c)
	db := store.MustNew(&c.DB)

	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Resource{},
		&model.UserRole{},
		&model.RoleResource{},
	); err != nil {
		log.Fatal().Err(err).Msg("auto migrate failed")
	}

	initDefaultData(db)

	if err := rbac.Init(db); err != nil {
		log.Fatal().Err(err).Msg("init casbin failed")
	}

	rbac.BatchRegisterRoutes(rbac.GetDefaultRouteMappings())
	middleware.SetTokenConfig(&c.Jwt)

	authService := service.NewAuth(db, ctx)
	authHandler := handler.NewAuth(authService)

	public := r.Group("/api/v1")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.Auth(), middleware.RBAC())
	{
		protected.GET("/users", func(c *gin.Context) {
			response.OK(c, gin.H{"message": "用户列表"})
		})
		protected.POST("/users", func(c *gin.Context) {
			response.OK(c, gin.H{"message": "创建用户"})
		})
	}

}

func initDefaultData(db *gorm.DB) {
	initResources(db)
	initRoles(db)
	initUsers(db)
}

func initResources(db *gorm.DB) {
	resources := []model.Resource{
		{Code: rbac.ResourceRole, Name: "角色管理", Actions: []string{rbac.ActionCreate, rbac.ActionRead, rbac.ActionUpdate, rbac.ActionDelete}},
		{Code: rbac.ResourceUser, Name: "用户管理", Actions: []string{rbac.ActionCreate, rbac.ActionRead, rbac.ActionUpdate, rbac.ActionDelete}},
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
