package dto

import (
	"time"

	"github.com/ningzining/cove/internal/pkg/model"
)

type ProfileResp struct {
	UserID    string              `json:"id"`
	Nickname  string              `json:"nickname"`
	Phone     string              `json:"phone"`
	Email     string              `json:"email"`
	Status    model.EnabledStatus `json:"status"`
	Roles     []RoleSimpleResp    `json:"roles"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

func ToProfileResp(user *model.User, roles []model.Role) *ProfileResp {
	roleList := make([]RoleSimpleResp, 0, len(roles))
	for _, r := range roles {
		roleList = append(roleList, RoleSimpleResp{
			RoleID: r.RoleID,
			Code:   r.Code,
			Name:   r.Name,
		})
	}
	return &ProfileResp{
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

type ProfileUpdateReq struct {
	Nickname string `json:"nickname" binding:"required,min=1,max=255"`
	Email    string `json:"email" binding:"omitempty,email,max=255"`
}

type ProfileUpdatePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=20"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}
