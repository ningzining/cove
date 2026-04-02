package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ningzining/cove/app/sys/internal/service/dto"
	"github.com/ningzining/cove/app/sys/internal/svc"
	"github.com/ningzining/cove/pkg/model"
	"github.com/ningzining/cove/pkg/token"
	"github.com/ningzining/cove/pkg/xerr"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Auth struct {
	DB  *gorm.DB
	ctx *svc.Context
}

func NewAuth(db *gorm.DB, ctx *svc.Context) *Auth {
	return &Auth{DB: db, ctx: ctx}
}

func (a *Auth) Login(req *dto.LoginReq) (*dto.LoginResp, error) {
	var user model.User
	err := a.DB.Model(&model.User{}).Where("phone = ?", req.Phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("phone", req.Phone).
				Str("provider", req.Provider).
				Msg("login: user not found")
			return nil, xerr.New(xerr.ErrLoginFailed)
		}
		log.Error().
			Err(err).
			Str("phone", req.Phone).
			Msg("login: db error")
		return nil, xerr.New(xerr.ErrDB)
	}
	// 校验密码
	if ok, err := a.compareHashAndPassword(user.Password, req.Password); !ok || err != nil {
		log.Error().
			Str("phone", req.Phone).
			Msg("login: invalid password")
		return nil, xerr.New(xerr.ErrLoginFailed)
	}
	// 校验用户状态
	if !user.Enabled() {
		log.Warn().
			Int64("user_id", user.ID).
			Str("phone", req.Phone).
			Msg("login: user disabled")
		return nil, xerr.New(xerr.ErrUserDisabled)
	}
	// 生成token
	now := time.Now()
	claims := token.CustomMapClaims{
		Provider: req.Provider,
		UserID:   user.ID,
		Phone:    user.Phone,
		Nickname: user.Nickname,
		Claims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(a.ctx.Config.Jwt.ExpireTime))),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    a.ctx.Config.Name,
		},
	}
	tokenString, err := token.Generate(claims, a.ctx.Config.Jwt.Key)
	if err != nil {
		log.Error().
			Err(err).
			Int64("user_id", user.ID).
			Msg("login: failed to generate token")
		return nil, xerr.New(xerr.ErrTokenSign)
	}

	return &dto.LoginResp{
		Token: tokenString,
	}, nil
}

func (a *Auth) Register(req *dto.RegisterReq) error {
	var existCount int64
	// 校验手机号是否已存在
	if err := a.DB.Model(&model.User{}).Where("phone = ?", req.Phone).Count(&existCount).Error; err != nil {
		log.Error().Err(err).Msg("Register err")
		return xerr.New(xerr.ErrDB)
	}
	if existCount > 0 {
		log.Error().Any("req", req).Msg("Register phone exist")
		return xerr.New(xerr.ErrPhoneExist)
	}
	user := req.Generate()
	if err := a.DB.Create(user).Error; err != nil {
		log.Error().Err(err).Msg("Register err")
		return xerr.New(xerr.ErrDB)
	}

	return nil
}

func (a *Auth) compareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}
