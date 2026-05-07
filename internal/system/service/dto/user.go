package dto

import (
	"time"

	"github.com/ningzining/cove/internal/pkg/model"
	"github.com/ningzining/cove/pkg/core/request"
)

type UserCreateReq struct {
	Nickname string   `json:"nickname" binding:"required,min=1,max=255"`
	Phone    string   `json:"phone" binding:"required,min=1,max=20"`
	Password string   `json:"password" binding:"required,min=6,max=20"`
	Email    string   `json:"email" binding:"omitempty,email,max=255"`
	RoleIDs  []string `json:"role_ids" binding:"omitempty"`
}

func (r *UserCreateReq) Generate() *model.User {
	return &model.User{
		Nickname: r.Nickname,
		Phone:    r.Phone,
		Password: r.Password,
		Email:    r.Email,
		Status:   model.Enabled,
	}
}

type UserDeleteReq struct {
	IDs []string `json:"ids" binding:"required"`
}

type UserUpdateReq struct {
	ID       string   `json:"-" binding:"omitempty"`
	Nickname string   `json:"nickname" binding:"required,min=1,max=255"`
	Phone    string   `json:"phone" binding:"required,min=1,max=20"`
	Email    string   `json:"email" binding:"omitempty,email,max=255"`
	RoleIDs  []string `json:"role_ids" binding:"omitempty"`
}

type UserPageReq struct {
	request.Pagination

	Nickname string `json:"nickname" form:"nickname" search:"type:icontains;column:nickname;table:sys_user"`
	Phone    string `json:"phone" form:"phone" search:"type:icontains;column:phone;table:sys_user"`
}

type UserUpdateStatusReq struct {
	ID     string              `json:"-" binding:"omitempty"`
	Status model.EnabledStatus `json:"status" binding:"required"`
}

type UserResetPasswordReq struct {
	ID          string `json:"-" binding:"omitempty"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}

type UserResp struct {
	UserID    string              `json:"id"`
	Nickname  string              `json:"nickname"`
	Phone     string              `json:"phone"`
	Email     string              `json:"email"`
	Status    model.EnabledStatus `json:"status"`
	Roles     []RoleSimpleResp    `json:"roles"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type RoleSimpleResp struct {
	RoleID string `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
}

func ToUserResp(user *model.User, roles []model.Role) *UserResp {
	roleList := make([]RoleSimpleResp, 0, len(roles))
	for _, r := range roles {
		roleList = append(roleList, RoleSimpleResp{
			RoleID: r.RoleID,
			Code:   r.Code,
			Name:   r.Name,
		})
	}
	return &UserResp{
		UserID:    user.UserID,
		Nickname:  user.Nickname,
		Phone:     user.Phone,
		Email:     user.Email,
		Status:    user.Status,
		Roles:     roleList,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
