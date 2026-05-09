package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/internal/pkg/middleware"
	"github.com/ningzining/cove/internal/pkg/response"
	"github.com/ningzining/cove/internal/pkg/xerr"
	"github.com/ningzining/cove/internal/system/service"
	"github.com/ningzining/cove/internal/system/service/dto"
	"github.com/rs/zerolog/log"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfile(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

// Get 获取个人信息
// @Summary 获取个人信息
// @Description 获取当前登录用户的个人信息
// @Tags 个人中心
// @Accept json
// @Produce json
// @Success 200 {object} response.response{data=dto.ProfileResp}
// @Router /api/v1/profile [get]
func (h *ProfileHandler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, xerr.New(xerr.ErrUnauthorized))
		return
	}

	profile, err := h.profileService.Get(userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, profile)
}

// Update 更新个人信息
// @Summary 更新个人信息
// @Description 更新当前登录用户的个人信息
// @Tags 个人中心
// @Accept json
// @Produce json
// @Param req body dto.ProfileUpdateReq true "更新个人信息请求"
// @Success 200 {object} response.response
// @Router /api/v1/profile [put]
func (h *ProfileHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, xerr.New(xerr.ErrUnauthorized))
		return
	}

	var req dto.ProfileUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}

	err := h.profileService.Update(userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// UpdatePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 个人中心
// @Accept json
// @Produce json
// @Param req body dto.ProfileUpdatePasswordReq true "修改密码请求"
// @Success 200 {object} response.response
// @Router /api/v1/profile/password [put]
func (h *ProfileHandler) UpdatePassword(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, xerr.New(xerr.ErrUnauthorized))
		return
	}

	var req dto.ProfileUpdatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("bind json failed")
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}

	err := h.profileService.UpdatePassword(userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
