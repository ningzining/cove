package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/internal/pkg/response"
	"github.com/ningzining/cove/internal/pkg/xerr"
	"github.com/ningzining/cove/internal/system/service"
	"github.com/ningzining/cove/internal/system/service/dto"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUser(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create 创建用户
// @Summary 创建用户
// @Description 创建用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param req body dto.UserCreateReq true "创建用户请求"
// @Success 200 {object} response.response
// @Router /api/v1/user [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	err := h.userService.Create(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 删除用户
// @Summary 删除用户
// @Description 删除用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param req body dto.UserDeleteReq true "删除用户请求"
// @Success 200 {object} response.response
// @Router /api/v1/user [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	var req dto.UserDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	err := h.userService.Delete(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Update 更新用户
// @Summary 更新用户
// @Description 更新用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param req body dto.UserUpdateReq true "更新用户请求"
// @Success 200 {object} response.response
// @Router /api/v1/user/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	var req dto.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	id := c.Param("id")
	req.ID = id
	err := h.userService.Update(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Get 获取用户详情
// @Summary 获取用户详情
// @Description 获取用户详情
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} response.response{data=dto.UserResp}
// @Router /api/v1/user/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.Get(id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, user)
}

// Page 分页查询用户
// @Summary 分页查询用户
// @Description 分页查询用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param req query dto.UserPageReq true "分页查询用户请求"
// @Success 200 {object} response.response{data=response.pageData{List=[]dto.UserResp}}
// @Router /api/v1/user [get]
func (h *UserHandler) Page(c *gin.Context) {
	var req dto.UserPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("bind query failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	users, total, err := h.userService.Page(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.PageOk(c, users, total)
}

// UpdateStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新用户状态
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param req body dto.UserUpdateStatusReq true "更新用户状态请求"
// @Success 200 {object} response.response
// @Router /api/v1/user/{id}/status [put]
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	var req dto.UserUpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	id := c.Param("id")
	req.ID = id
	err := h.userService.UpdateStatus(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// ResetPassword 重置用户密码
// @Summary 重置用户密码
// @Description 重置用户密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param req body dto.UserResetPasswordReq true "重置密码请求"
// @Success 200 {object} response.response
// @Router /api/v1/user/{id}/password/reset [put]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req dto.UserResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	id := c.Param("id")
	req.ID = id
	err := h.userService.ResetPassword(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
