package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type EnabledStatus int8

const (
	Enabled  EnabledStatus = 1
	Disabled EnabledStatus = 2
)

type User struct {
	ID        int64          `json:"-" gorm:"primaryKey"`
	UserID    string         `json:"user_id" gorm:"varchar(36);unique;not null;comment:用户ID"`
	Nickname  string         `json:"nickname" gorm:"varchar(255);not null;comment:昵称"`
	Phone     string         `json:"phone" gorm:"varchar(20);unique;not null;comment:手机号"`
	Password  string         `json:"-" gorm:"varchar(255);not null;comment:密码"`
	Email     string         `json:"email" gorm:"varchar(255);comment:邮箱"`
	Status    EnabledStatus  `json:"status" gorm:"not null;default:1;comment:状态,1:启用,2:禁用"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

func (u *User) Enabled() bool {
	return u.Status == Enabled
}

func (u *User) TableName() string {
	return "sys_user"
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	// 生成用户ID
	uid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.UserID = uid.String()

	// 加密密码
	if err := u.Encrypt(); err != nil {
		return err
	}

	return nil
}

func (u *User) BeforeUpdate(_ *gorm.DB) error {
	// 加密密码
	if err := u.Encrypt(); err != nil {
		return err
	}
	return nil
}

func (u *User) Encrypt() error {
	if u.Password == "" {
		return nil
	}

	// 加密密码
	var hash []byte
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
