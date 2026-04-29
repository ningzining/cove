package service

import (
	"errors"

	"github.com/ningzining/cove/internal/system/service/dto"
	"github.com/ningzining/cove/internal/system/svc"
	"github.com/ningzining/cove/pkg/core/casbin"
	"github.com/ningzining/cove/pkg/core/search"
	"github.com/ningzining/cove/pkg/model"
	"github.com/ningzining/cove/pkg/xerr"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type RoleService struct {
	DB  *gorm.DB
	ctx *svc.Context
}

func NewRole(db *gorm.DB, ctx *svc.Context) *RoleService {
	return &RoleService{DB: db, ctx: ctx}
}

func (s *RoleService) Create(req *dto.RoleCreateReq) error {
	var err error
	tx := s.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 校验角色编码是否存在
	{
		var count int64
		err = tx.Model((*model.Role)(nil)).Where("code = ?", req.Code).Count(&count).Error
		if err != nil {
			log.Err(err).Str("code", req.Code).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}
		if count > 0 {
			return xerr.New(xerr.ErrRoleCodeExist)
		}
	}
	// 校验角色名称是否存在
	{
		var count int64
		err = tx.Model((*model.Role)(nil)).Where("name = ?", req.Name).Count(&count).Error
		if err != nil {
			log.Err(err).Str("name", req.Name).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}
		if count > 0 {
			return xerr.New(xerr.ErrRoleNameExist)
		}
	}
	// 查询资源是否存在
	var resources []model.Resource
	if err = tx.Where("resource_id IN ?", req.ResourceIDs).Find(&resources).Error; err != nil {
		log.Err(err).Any("resource_ids", req.ResourceIDs).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}
	// 创建角色
	role := req.Generate()
	if err = tx.Create(role).Error; err != nil {
		log.Err(err).Any("role", role).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}
	// 创建角色资源关联
	roleResources := make([]model.RoleResource, 0)
	polices := make([][]string, 0)
	for _, v := range resources {
		roleResources = append(roleResources, model.RoleResource{
			RoleID:     role.RoleID,
			ResourceID: v.ResourceID,
		})
		polices = append(polices, []string{role.Code, v.Code, string(v.Action)})
	}
	if len(roleResources) > 0 {
		if err = tx.Create(&roleResources).Error; err != nil {
			log.Err(err).Any("roleResources", roleResources).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}
	}
	// 保存策略
	if len(polices) > 0 {
		_, err = casbin.Enforcer().AddNamedPolicies("p", polices)
		if err != nil {
			log.Err(err).
				Str("role_code", role.Code).
				Any("policies", polices).
				Msg("add casbin policy failed")
			return xerr.New(xerr.ErrDB)
		}
	}

	return nil
}

func (s *RoleService) Delete(req *dto.RoleDeleteReq) error {
	if len(req.IDs) == 0 {
		return nil
	}
	var err error
	tx := s.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	roles := make([]model.Role, 0)
	if err = tx.Where("role_id IN ?", req.IDs).Find(&roles).Error; err != nil {
		log.Err(err).Any("ids", req.IDs).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}
	// 内置角色不能操作
	for _, v := range roles {
		if v.Code == model.NormalUserRoleCode || v.Code == model.AdminRoleCode {
			return xerr.New(xerr.ErrRoleCannotOperate)
		}
	}

	// 检查是否有用户正在使用这些角色
	var count int64
	err = tx.Model((*model.UserRole)(nil)).Where("role_id IN ?", req.IDs).Count(&count).Error
	if err != nil {
		log.Err(err).Any("ids", req.IDs).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}
	if count > 0 {
		return xerr.New(xerr.ErrRoleInUse)
	}
	// 删除角色
	err = tx.Delete(&model.Role{}, "role_id IN ?", req.IDs).Error
	if err != nil {
		log.Err(err).Any("ids", req.IDs).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}
	// 删除 Casbin 策略
	for _, role := range roles {
		_, err = casbin.Enforcer().RemoveFilteredPolicy(0, role.Code)
		if err != nil {
			log.Err(err).Str("role_code", role.Code).Msg("remove casbin policy failed")
			return xerr.New(xerr.ErrDB)
		}
	}

	return nil
}

func (s *RoleService) Update(req *dto.RoleUpdateReq) error {
	var err error
	tx := s.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 校验角色是否存在
	var role model.Role
	err = tx.Where("role_id = ?", req.ID).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xerr.New(xerr.ErrRoleNotExist)
		}
		log.Err(err).Str("id", req.ID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}
	// 内置角色不能操作
	if role.Code == model.NormalUserRoleCode || role.Code == model.AdminRoleCode {
		return xerr.New(xerr.ErrRoleCannotOperate)
	}

	// 如果名称有变化，检查是否与其他角色重复
	if req.Name != role.Name {
		var count int64
		err = tx.Model((*model.Role)(nil)).Where("name = ? AND role_id != ?", req.Name, req.ID).Count(&count).Error
		if err != nil {
			log.Err(err).Str("name", req.Name).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}
		if count > 0 {
			return xerr.New(xerr.ErrRoleNameExist)
		}
	}

	// 更新角色
	role.Name = req.Name
	if err = tx.Save(&role).Error; err != nil {
		return xerr.New(xerr.ErrDB)
	}

	// 查询资源是否存在
	var resources []model.Resource
	if err = tx.Where("resource_id IN ?", req.ResourceIDs).Find(&resources).Error; err != nil {
		log.Err(err).Any("resource_ids", req.ResourceIDs).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}
	// 删除角色资源关联
	err = tx.Delete(&model.RoleResource{}, "role_id = ?", role.RoleID).Error
	if err != nil {
		log.Err(err).Str("role_id", role.RoleID).Msg("delete role resource failed")
		return xerr.New(xerr.ErrDB)
	}
	// 创建角色资源关联
	roleResources := make([]model.RoleResource, 0)
	polices := make([][]string, 0)
	for _, v := range resources {
		roleResources = append(roleResources, model.RoleResource{
			RoleID:     role.RoleID,
			ResourceID: v.ResourceID,
		})
		polices = append(polices, []string{role.Code, v.Code, string(v.Action)})
	}
	if len(roleResources) > 0 {
		if err = tx.Create(&roleResources).Error; err != nil {
			log.Err(err).Any("roleResources", roleResources).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}
	}

	// 更新 Casbin 策略：先删后加
	_, err = casbin.Enforcer().RemoveFilteredPolicy(0, role.Code)
	if err != nil {
		log.Err(err).Str("role_code", role.Code).Msg("remove casbin policy failed")
		return xerr.New(xerr.ErrDB)
	}

	// 批量添加策略
	if len(polices) > 0 {
		_, err = casbin.Enforcer().AddNamedPolicies("p", polices)
		if err != nil {
			log.Err(err).
				Str("role_id", role.RoleID).
				Str("role_code", role.Code).
				Any("policies", polices).
				Msg("add casbin policy failed")
			return xerr.New(xerr.ErrDB)
		}
	}

	return nil
}

func (s *RoleService) Page(req *dto.RolePageReq) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64
	err := s.DB.Model((*model.Role)(nil)).
		Scopes(search.MakeCondition(req)).
		Count(&total).Error
	if err != nil {
		log.Err(err).Any("req", req).Msg("db error")
		return nil, 0, xerr.New(xerr.ErrDB)
	}

	err = s.DB.Model((*model.Role)(nil)).
		Scopes(
			search.MakeCondition(req),
			search.Paginate(req.GetPage(), req.GetPageSize()),
		).
		Order("id DESC").
		Find(&roles).Error
	if err != nil {
		log.Err(err).Any("req", req).Msg("db error")
		return nil, 0, xerr.New(xerr.ErrDB)
	}
	return roles, total, nil
}

// RoleWithPermissions 带权限信息的角色
type RoleWithPermissions struct {
	*model.Role
	Permissions []RolePermission `json:"permissions,omitempty"`
}

// RolePermission 角色权限
type RolePermission struct {
	ResourceID string   `json:"resource_id"`
	Code       string   `json:"code"`
	Name       string   `json:"name"`
	Actions    []string `json:"actions"`
}

func (s *RoleService) Get(id string) (*RoleWithPermissions, error) {
	var role model.Role
	err := s.DB.Where("role_id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.New(xerr.ErrRoleNotExist)
		}
		log.Err(err).Str("id", id).Msg("db error")
		return nil, xerr.New(xerr.ErrDB)
	}

	result := &RoleWithPermissions{Role: &role}

	// 从 Casbin 读取策略
	enforcer := casbin.Enforcer()
	if enforcer != nil {
		policies, _ := enforcer.GetFilteredPolicy(0, id)
		if len(policies) > 0 {
			// 按 resource code 分组
			resourceActions := make(map[string][]string)
			for _, p := range policies {
				if len(p) >= 3 {
					resourceCode := p[1]
					action := p[2]
					resourceActions[resourceCode] = append(resourceActions[resourceCode], action)
				}
			}

			// 查询 Resource 表获取详情
			if len(resourceActions) > 0 {
				codes := make([]string, 0, len(resourceActions))
				for code := range resourceActions {
					codes = append(codes, code)
				}

				var resources []*model.Resource
				err = s.DB.Where("code IN ?", codes).Find(&resources).Error
				if err != nil {
					log.Err(err).Any("codes", codes).Msg("query resources failed")
				} else {
					for _, res := range resources {
						if actions, ok := resourceActions[res.Code]; ok {
							result.Permissions = append(result.Permissions, RolePermission{
								ResourceID: res.ResourceID,
								Code:       res.Code,
								Name:       res.Name,
								Actions:    actions,
							})
						}
					}
				}
			}
		}
	}

	return result, nil
}

func (s *RoleService) Option() ([]*model.Role, error) {
	var roles []*model.Role
	err := s.DB.Order("id DESC").Find(&roles).Error
	if err != nil {
		log.Err(err).Msg("db error")
		return nil, xerr.New(xerr.ErrDB)
	}
	return roles, nil
}

func (s *RoleService) UpdateStatus(req *dto.RoleUpdateStatusReq) error {
	var err error
	tx := s.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var role model.Role
	err = tx.Where("role_id = ?", req.ID).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xerr.New(xerr.ErrRoleNotExist)
		}
		log.Err(err).Str("role_id", req.ID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}
	// 内置角色不能操作
	if role.Code == model.NormalUserRoleCode || role.Code == model.AdminRoleCode {
		return xerr.New(xerr.ErrRoleCannotOperate)
	}
	// 状态未改变，无需更新
	if req.Status == role.Status {
		return nil
	}
	// 更新角色状态
	role.Status = req.Status
	if err = tx.Save(&role).Error; err != nil {
		log.Err(err).Any("status", req.Status).Str("role_id", req.ID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	if req.Status == model.Enabled {
		// 同步角色权限到 Casbin
		resources := make([]model.Resource, 0)
		err = tx.Where("resource_id in (select resource_id from sys_role_resource where role_id = ?)", role.RoleID).Find(&resources).Error
		if err != nil {
			log.Err(err).Str("role_id", role.RoleID).Msg("query resources failed")
			return xerr.New(xerr.ErrDB)
		}
		policies := make([][]string, 0, len(resources))
		for _, v := range resources {
			policies = append(policies, []string{role.Code, v.Code, string(v.Action)})
		}
		// 加载策略
		_, err = casbin.Enforcer().AddNamedPolicies("p", policies)
		if err != nil {
			return err
		}
	} else {
		// 移除角色权限从 Casbin
		_, err = casbin.Enforcer().RemoveFilteredPolicy(0, role.Code)
		if err != nil {
			return err
		}
	}

	return nil
}
