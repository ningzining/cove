package service

import (
	"errors"

	"github.com/ningzining/cove/internal/pkg/model"
	"github.com/ningzining/cove/internal/pkg/xerr"
	"github.com/ningzining/cove/internal/system/service/dto"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ProfileService struct {
	db *gorm.DB
}

func NewProfile(db *gorm.DB) *ProfileService {
	return &ProfileService{db: db}
}

func (s *ProfileService) Get(userID string) (*dto.ProfileResp, error) {
	var user model.User
	err := s.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.New(xerr.ErrUserNotExist)
		}
		log.Error().Err(err).Str("user_id", userID).Msg("db error")
		return nil, xerr.New(xerr.ErrDB)
	}

	var userRoles []model.UserRole
	err = s.db.Where("user_id = ?", user.UserID).Find(&userRoles).Error
	if err != nil {
		log.Error().Err(err).Str("user_id", user.UserID).Msg("db error")
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
			log.Error().Err(err).Any("role_ids", roleIDs).Msg("db error")
			return nil, xerr.New(xerr.ErrDB)
		}
	}

	return dto.ToProfileResp(&user, roles), nil
}

func (s *ProfileService) Update(userID string, req *dto.ProfileUpdateReq) error {
	var user model.User
	err := s.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xerr.New(xerr.ErrUserNotExist)
		}
		log.Error().Err(err).Str("user_id", userID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	user.Nickname = req.Nickname
	user.Email = req.Email

	if err := s.db.Omit("password").Save(&user).Error; err != nil {
		log.Error().Err(err).Any("user", user).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	return nil
}

func (s *ProfileService) UpdatePassword(userID string, req *dto.ProfileUpdatePasswordReq) error {
	var user model.User
	err := s.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xerr.New(xerr.ErrUserNotExist)
		}
		log.Error().Err(err).Str("user_id", userID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	if ok, _ := s.compareHashAndPassword(user.Password, req.NewPassword); ok {
		return xerr.New(xerr.ErrNewPasswordSame)
	}

	if ok, _ := s.compareHashAndPassword(user.Password, req.OldPassword); !ok {
		return xerr.New(xerr.ErrOldPasswordIncorrect)
	}

	user.Password = req.NewPassword
	if err := s.db.Save(&user).Error; err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("db error")
		return xerr.New(xerr.ErrDB)
	}

	return nil
}

func (s *ProfileService) compareHashAndPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	log.Info().Err(err).Str("hashed_password", hashedPassword).Str("password", password).Msg("compare hash and password")

	if err != nil {
		return false, err
	}
	log.Info().Bool("matched", true).Str("hashed_password", hashedPassword).Str("password", password).Msg("compare hash and password")
	return true, nil
}
