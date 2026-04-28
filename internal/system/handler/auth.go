package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/internal/system/service"
	"github.com/ningzining/cove/internal/system/service/dto"
	"github.com/ningzining/cove/pkg/rest/response"
	"github.com/ningzining/cove/pkg/xerr"
)

type Auth struct {
	authService *service.Auth
}

func NewAuth(authService *service.Auth) *Auth {
	return &Auth{authService: authService}
}

// Login 登录
// @Summary 登录
// @Description 登录用户
// @Tags Auth
// @Accept json
// @Produce json
// @Param req body dto.LoginReq true "登录请求"
// @Success 200 {object} response.response{data=dto.LoginResp}
// @Router /api/v1/login [post]
func (a *Auth) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, xerr.New(xerr.ErrBind))
		return
	}
	data, err := a.authService.Login(&req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}
