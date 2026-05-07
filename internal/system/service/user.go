package service

import (
	"errors"

	"github.com/ningzining/cove/internal/pkg/model"
	"github.com/ningzining/cove/internal/pkg/xerr"
	"github.com/ningzining/cove/internal/system/service/dto"
	"github.com/ningzining/cove/internal/system/svc"
	"github.com/ningzining/cove/pkg/core/casbin"
	"github.com/ningzining/cove/pkg/core/search"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserService struct {
	db  *gorm.DB
	ctx *svc.Context
}

func NewUser(db *gorm.DB, ctx *svc.Context) *UserService {
	return &UserService{db: db, ctx: ctx}
}

func (s *UserService) Create(req *dto.UserCreateReq) error {
	var err error
	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 校验手机号是否存在
	var count int64
	err = tx.Model((*model.User)(nil)).Where("phone = ?", req.Phone).Count(&count).Error
	if err != nil {
		log.Err(err).Str("phone", req.Phone).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}
	if count > 0 {
		return xerr.New(xerr.ErrUserPhoneExist)
	}

	// 校验角色是否存在
	var roles []model.Role
	if len(req.RoleIDs) > 0 {
		err = tx.Where("role_id IN ?", req.RoleIDs).Find(&roles).Error
		if err != nil {
			log.Err(err).Any("role_ids", req.RoleIDs).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}
		if len(roles) != len(req.RoleIDs) {
			return xerr.New(xerr.ErrRoleNotExist)
		}
	}

	// 创建用户
	user := req.Generate()
	if err = tx.Create(user).Error; err != nil {
		log.Err(err).Any("user", user).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	// 创建用户-角色关联
	if len(roles) > 0 {
		userRoles := make([]model.UserRole, 0, len(roles))
		for _, r := range roles {
			userRoles = append(userRoles, model.UserRole{
				UserID: user.UserID,
				RoleID: r.RoleID,
			})
		}
		if err = tx.Create(&userRoles).Error; err != nil {
			log.Err(err).Any("user_roles", userRoles).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}

		// 同步到 Casbin
		roleCodes := make([]string, 0, len(roles))
		for _, r := range roles {
			roleCodes = append(roleCodes, r.Code)
		}
		_, err = casbin.Enforcer().AddRolesForUser(user.UserID, roleCodes)
		if err != nil {
			log.Err(err).
				Str("user_id", user.UserID).
				Any("role_codes", roleCodes).
				Msg("add role for user failed")
			return xerr.New(xerr.ErrDB)
		}
	}

	return nil
}

func (s *UserService) Delete(req *dto.UserDeleteReq) error {
	if len(req.IDs) == 0 {
		return nil
	}
	var err error
	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 查询用户
	var users []model.User
	err = tx.Where("user_id IN ?", req.IDs).Find(&users).Error
	if err != nil {
		log.Err(err).Any("ids", req.IDs).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	// 校验内置用户不能删除
	for _, u := range users {
		if u.Phone == model.AdminPhone || u.Phone == model.NormalUserPhone {
			return xerr.New(xerr.ErrUserCannotOperate)
		}
	}

	// 查询用户关联的角色
	userIDs := make([]string, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.UserID)
	}
	var userRoles []model.UserRole
	err = tx.Where("user_id IN ?", userIDs).Find(&userRoles).Error
	if err != nil {
		log.Err(err).Any("user_ids", userIDs).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	// 查询角色信息
	roleIDs := make([]string, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}
	var roles []model.Role
	if len(roleIDs) > 0 {
		err = tx.Where("role_id IN ?", roleIDs).Find(&roles).Error
		if err != nil {
			log.Err(err).Any("role_ids", roleIDs).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}
	}
	roleMap := make(map[string]model.Role)
	for _, r := range roles {
		roleMap[r.RoleID] = r
	}

	// 删除用户-角色关联
	err = tx.Delete(&model.UserRole{}, "user_id IN ?", userIDs).Error
	if err != nil {
		log.Err(err).Any("user_ids", userIDs).Msg("delete user role failed")
		return xerr.New(xerr.ErrDB)
	}

	// 逻辑删除用户
	err = tx.Delete(&model.User{}, "user_id IN ?", userIDs).Error
	if err != nil {
		log.Err(err).Any("user_ids", userIDs).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	// 移除 Casbin 中的用户角色分配
	for _, ur := range userRoles {
		_, err = casbin.Enforcer().DeleteRolesForUser(ur.UserID)
		if err != nil {
			log.Err(err).
				Str("user_id", ur.UserID).
				Msg("delete role for user failed")
			return xerr.New(xerr.ErrDB)
		}
	}

	return nil
}

func (s *UserService) Update(req *dto.UserUpdateReq) error {
	var err error
	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 校验用户是否存在
	var user model.User
	err = tx.Where("user_id = ?", req.ID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xerr.New(xerr.ErrUserNotExist)
		}
		log.Err(err).Str("id", req.ID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	// 校验内置用户不能修改
	if user.Phone == model.AdminPhone || user.Phone == model.NormalUserPhone {
		return xerr.New(xerr.ErrUserCannotOperate)
	}

	// 校验手机号是否与其他用户重复
	if req.Phone != user.Phone {
		var count int64
		err = tx.Model((*model.User)(nil)).Where("phone = ? AND user_id != ?", req.Phone, req.ID).Count(&count).Error
		if err != nil {
			log.Err(err).Str("phone", req.Phone).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}
		if count > 0 {
			return xerr.New(xerr.ErrUserPhoneExist)
		}
	}

	// 更新用户基本信息
	user.Nickname = req.Nickname
	user.Phone = req.Phone
	user.Email = req.Email
	if err = tx.Save(&user).Error; err != nil {
		log.Err(err).Any("user", user).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	// 如果传了 role_ids，则更新角色
	if req.RoleIDs != nil {
		// 查询新角色是否存在
		var newRoles []model.Role
		if len(req.RoleIDs) > 0 {
			err = tx.Where("role_id IN ?", req.RoleIDs).Find(&newRoles).Error
			if err != nil {
				log.Err(err).Any("role_ids", req.RoleIDs).Msg("db error")
				return xerr.New(xerr.ErrDB)
			}
			if len(newRoles) != len(req.RoleIDs) {
				return xerr.New(xerr.ErrRoleNotExist)
			}
		}

		// 查询旧角色
		var oldUserRoles []model.UserRole
		err = tx.Where("user_id = ?", user.UserID).Find(&oldUserRoles).Error
		if err != nil {
			log.Err(err).Str("user_id", user.UserID).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}
		var oldRoles []model.Role
		if len(oldUserRoles) > 0 {
			oldRoleIDs := make([]string, 0, len(oldUserRoles))
			for _, ur := range oldUserRoles {
				oldRoleIDs = append(oldRoleIDs, ur.RoleID)
			}
			err = tx.Where("role_id IN ?", oldRoleIDs).Find(&oldRoles).Error
			if err != nil {
				log.Err(err).Any("old_role_ids", oldRoleIDs).Msg("db error")
				return xerr.New(xerr.ErrDB)
			}
		}

		// 删除旧关联
		err = tx.Delete(&model.UserRole{}, "user_id = ?", user.UserID).Error
		if err != nil {
			log.Err(err).Str("user_id", user.UserID).Msg("delete user role failed")
			return xerr.New(xerr.ErrDB)
		}

		// 创建新关联
		if len(newRoles) > 0 {
			userRoles := make([]model.UserRole, 0, len(newRoles))
			for _, r := range newRoles {
				userRoles = append(userRoles, model.UserRole{
					UserID: user.UserID,
					RoleID: r.RoleID,
				})
			}
			if err = tx.Create(&userRoles).Error; err != nil {
				log.Err(err).Any("user_roles", userRoles).Msg("db error")
				return xerr.New(xerr.ErrDB)
			}
		}

		// 同步 Casbin：先删旧的，再加新的
		_, err = casbin.Enforcer().DeleteRolesForUser(user.UserID)
		if err != nil {
			log.Err(err).
				Str("user_id", user.UserID).
				Msg("delete role for user failed")
			return xerr.New(xerr.ErrDB)
		}

		roleCodes := make([]string, 0, len(newRoles))
		for _, r := range newRoles {
			roleCodes = append(roleCodes, r.Code)
		}
		if len(roleCodes) > 0 {
			_, err = casbin.Enforcer().AddRolesForUser(user.UserID, roleCodes)
			if err != nil {
				log.Err(err).
					Str("user_id", user.UserID).
					Any("role_codes", roleCodes).
					Msg("add role for user failed")
				return xerr.New(xerr.ErrDB)
			}
		}
	}

	return nil
}

func (s *UserService) Page(req *dto.UserPageReq) ([]*dto.UserResp, int64, error) {
	var users []*model.User
	var total int64
	err := s.db.Model((*model.User)(nil)).
		Scopes(search.MakeCondition(req)).
		Count(&total).Error
	if err != nil {
		log.Err(err).Any("req", req).Msg("db error")
		return nil, 0, xerr.New(xerr.ErrDB)
	}

	err = s.db.Model((*model.User)(nil)).
		Scopes(
			search.MakeCondition(req),
			search.Paginate(req.GetPage(), req.GetPageSize()),
		).
		Order("created_at DESC, id DESC").
		Find(&users).Error
	if err != nil {
		log.Err(err).Any("req", req).Msg("db error")
		return nil, 0, xerr.New(xerr.ErrDB)
	}

	// 加载用户的角色
	userIDs := make([]string, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.UserID)
	}
	var userRoles []model.UserRole
	if len(userIDs) > 0 {
		err = s.db.Where("user_id IN ?", userIDs).Find(&userRoles).Error
		if err != nil {
			log.Err(err).Any("user_ids", userIDs).Msg("db error")
			return nil, 0, xerr.New(xerr.ErrDB)
		}
	}
	roleIDs := make([]string, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}
	var roles []model.Role
	if len(roleIDs) > 0 {
		err = s.db.Where("role_id IN ?", roleIDs).Find(&roles).Error
		if err != nil {
			log.Err(err).Any("role_ids", roleIDs).Msg("db error")
			return nil, 0, xerr.New(xerr.ErrDB)
		}
	}

	// 组装数据
	userRoleMap := make(map[string][]model.Role)
	roleMap := make(map[string]model.Role)
	for _, r := range roles {
		roleMap[r.RoleID] = r
	}
	for _, ur := range userRoles {
		if role, ok := roleMap[ur.RoleID]; ok {
			userRoleMap[ur.UserID] = append(userRoleMap[ur.UserID], role)
		}
	}

	result := make([]*dto.UserResp, 0, len(users))
	for _, u := range users {
		userRoles := userRoleMap[u.UserID]
		result = append(result, dto.ToUserResp(u, userRoles))
	}

	return result, total, nil
}

func (s *UserService) Get(id string) (*dto.UserResp, error) {
	var user model.User
	err := s.db.Where("user_id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.New(xerr.ErrUserNotExist)
		}
		log.Err(err).Str("id", id).Msg("db error")
		return nil, xerr.New(xerr.ErrDB)
	}

	// 加载用户的角色
	var userRoles []model.UserRole
	err = s.db.Where("user_id = ?", user.UserID).Find(&userRoles).Error
	if err != nil {
		log.Err(err).Str("user_id", user.UserID).Msg("db error")
		return nil, xerr.New(xerr.ErrDB)
	}
	var roles []model.Role
	if len(userRoles) > 0 {
		roleIDs := make([]string, 0, len(userRoles))
		for _, ur := range userRoles {
			roleIDs = append(roleIDs, ur.RoleID)
		}
		err = s.db.Where("role_id IN ?", roleIDs).Find(&roles).Error
		if err != nil {
			log.Err(err).Any("role_ids", roleIDs).Msg("db error")
			return nil, xerr.New(xerr.ErrDB)
		}
	}

	return dto.ToUserResp(&user, roles), nil
}

func (s *UserService) UpdateStatus(req *dto.UserUpdateStatusReq) error {
	var err error
	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var user model.User
	err = tx.Where("user_id = ?", req.ID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xerr.New(xerr.ErrUserNotExist)
		}
		log.Err(err).Str("user_id", req.ID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	// 校验内置用户不能操作
	if user.Phone == model.AdminPhone || user.Phone == model.NormalUserPhone {
		return xerr.New(xerr.ErrUserCannotOperate)
	}

	// 状态未改变，无需更新
	if req.Status == user.Status {
		return nil
	}

	// 更新用户状态
	user.Status = req.Status
	if err = tx.Save(&user).Error; err != nil {
		log.Err(err).Any("status", req.Status).Str("user_id", req.ID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	// 查询用户关联的角色
	var userRoles []model.UserRole
	err = tx.Where("user_id = ?", user.UserID).Find(&userRoles).Error
	if err != nil {
		log.Err(err).Str("user_id", user.UserID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}
	var roles []model.Role
	if len(userRoles) > 0 {
		roleIDs := make([]string, 0, len(userRoles))
		for _, ur := range userRoles {
			roleIDs = append(roleIDs, ur.RoleID)
		}
		err = tx.Where("role_id IN ?", roleIDs).Find(&roles).Error
		if err != nil {
			log.Err(err).Any("role_ids", roleIDs).Msg("db error")
			return xerr.New(xerr.ErrDB)
		}
	}

	if req.Status == model.Enabled {
		// 启用用户：恢复 Casbin 角色分配
		roleCodes := make([]string, 0, len(roles))
		for _, r := range roles {
			roleCodes = append(roleCodes, r.Code)
		}
		_, err = casbin.Enforcer().AddRolesForUser(user.UserID, roleCodes)
		if err != nil {
			log.Err(err).
				Str("user_id", user.UserID).
				Any("role_codes", roleCodes).
				Msg("add role for user failed")
			return xerr.New(xerr.ErrDB)
		}
	} else {
		// 禁用用户：移除 Casbin 角色分配
		_, err = casbin.Enforcer().DeleteRolesForUser(user.UserID)
		if err != nil {
			log.Err(err).
				Str("user_id", user.UserID).
				Msg("delete role for user failed")
			return xerr.New(xerr.ErrDB)
		}
	}

	return nil
}

func (s *UserService) ResetPassword(req *dto.UserResetPasswordReq) error {
	var err error
	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var user model.User
	err = tx.Where("user_id = ?", req.ID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xerr.New(xerr.ErrUserNotExist)
		}
		log.Err(err).Str("id", req.ID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	// 校验内置用户不能操作
	if user.Phone == model.AdminPhone || user.Phone == model.NormalUserPhone {
		return xerr.New(xerr.ErrUserCannotOperate)
	}

	user.Password = req.NewPassword
	if err = tx.Save(&user).Error; err != nil {
		log.Err(err).Str("user_id", req.ID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	return nil
}
