package dto

import (
	"github.com/ningzining/cove/pkg/model"
)

type RegisterReq struct {
	Nickname string `json:"nickname" binding:"required,mix=6,max=20"`
	Phone    string `json:"phone" binding:"required,mix=6,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
}

func (r *RegisterReq) Generate() *model.User {
	return &model.User{
		Nickname: r.Nickname,
		Phone:    r.Phone,
		Password: r.Password,
		Email:    r.Email,
	}
}

type LoginReq struct {
	Provider string `json:"-"` // 登录方式

	Phone    string `json:"phone" binding:"required,mix=6,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type LoginResp struct {
	Token string `json:"token"` // 登录token
}
