package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/internal/system/service"
	"github.com/ningzining/cove/internal/system/service/dto"
	"github.com/ningzining/cove/pkg/rest/response"
	"github.com/ningzining/cove/pkg/xerr"
	"github.com/rs/zerolog/log"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRole(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// Create 创建角色
// @Summary 创建角色
// @Description 创建角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param req body dto.RoleCreateReq true "创建角色请求"
// @Success 200 {object} response.response
// @Router /api/v1/role [post]
func (r *RoleHandler) Create(c *gin.Context) {
	var req dto.RoleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	err := r.roleService.Create(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 删除角色
// @Summary 删除角色
// @Description 删除角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param req body dto.RoleDeleteReq true "删除角色请求"
// @Success 200 {object} response.response
// @Router /api/v1/role [delete]
func (r *RoleHandler) Delete(c *gin.Context) {
	var req dto.RoleDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	err := r.roleService.Delete(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Update 更新角色
// @Summary 更新角色
// @Description 更新角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Param req body dto.RoleUpdateReq true "更新角色请求"
// @Success 200 {object} response.response
// @Router /api/v1/role/{id} [put]
func (r *RoleHandler) Update(c *gin.Context) {
	var req dto.RoleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	id := c.Param("id")
	req.ID = id
	err := r.roleService.Update(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Get 获取角色详情
// @Summary 获取角色详情
// @Description 获取角色详情
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Success 200 {object} response.response
// @Router /api/v1/role/{id} [get]
func (r *RoleHandler) Get(c *gin.Context) {
	id := c.Param("id")
	role, err := r.roleService.Get(id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, role)
}

// Page 分页查询角色
// @Summary 分页查询角色
// @Description 分页查询角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param req query dto.RolePageReq true "分页查询角色请求"
// @Success 200 {object} response.response{data=response.pageData{List=[]model.Role}}
// @Router /api/v1/role [get]
func (r *RoleHandler) Page(c *gin.Context) {
	var req dto.RolePageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("bind query failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	roles, total, err := r.roleService.Page(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.PageOk(c, roles, total)
}

// Option 角色选项
// @Summary 角色选项
// @Description 角色选项
// @Tags 角色
// @Accept json
// @Produce json
// @Success 200 {object} response.response{data=[]model.Role}
// @Router /api/v1/role/option [get]
func (r *RoleHandler) Option(c *gin.Context) {
	roles, err := r.roleService.Option()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, roles)
}

// UpdateStatus 更新角色状态
// @Summary 更新角色状态
// @Description 更新角色状态
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Param req body dto.RoleUpdateStatusReq true "更新角色状态请求"
// @Success 200 {object} response.response
// @Router /api/v1/role/{id}/status [put]
func (r *RoleHandler) UpdateStatus(c *gin.Context) {
	var req dto.RoleUpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	id := c.Param("id")
	req.ID = id
	err := r.roleService.UpdateStatus(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
