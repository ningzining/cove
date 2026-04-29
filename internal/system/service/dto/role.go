package dto

import (
	"github.com/ningzining/cove/pkg/core/request"
	"github.com/ningzining/cove/pkg/model"
)

type RoleCreateReq struct {
	Code        string   `json:"code" binding:"required,min=1,max=64"`
	Name        string   `json:"name" binding:"required,min=1,max=255"`
	ResourceIDs []string `json:"resource_ids" binding:"omitempty"`
}

func (r *RoleCreateReq) Generate() *model.Role {
	return &model.Role{
		Code:   r.Code,
		Name:   r.Name,
		Status: model.Enabled,
	}
}

type RoleDeleteReq struct {
	IDs []string `json:"ids" binding:"required"`
}

type RoleUpdateReq struct {
	ID          string   `json:"-" binding:"omitempty"`
	Name        string   `json:"name" binding:"required,min=1,max=255"`
	ResourceIDs []string `json:"resource_ids" binding:"omitempty"`
}

type RolePageReq struct {
	request.Pagination

	Name string `json:"name" form:"name" search:"type:icontains;column:name;table:sys_role"`
}

type RoleUpdateStatusReq struct {
	ID     string              `json:"-" binding:"omitempty"`
	Status model.EnabledStatus `json:"status" binding:"omitempty"`
}
