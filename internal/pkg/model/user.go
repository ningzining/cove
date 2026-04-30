package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	AdminNickname = "管理员"         // 管理员昵称
	AdminPhone    = "13800000000" // 管理员手机号

	NormalUserNickname = "普通用户"        // 普通用户昵称
	NormalUserPhone    = "13800000001" // 普通用户手机号

	DefaultPassword = "Cove@123456" // 默认密码
)

type EnabledStatus int8

const (
	Enabled  EnabledStatus = 1 // 启用
	Disabled EnabledStatus = 2 // 禁用
)

type User struct {
	ID        int64          `json:"-" gorm:"autoIncrement;primaryKey"`
	UserID    string         `json:"id" gorm:"type:varchar(64);unique;not null;comment:用户ID"`
	Nickname  string         `json:"nickname" gorm:"type:varchar(255);not null;comment:昵称"`
	Phone     string         `json:"phone" gorm:"type:varchar(20);not null;comment:手机号"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null;comment:密码"`
	Email     string         `json:"email" gorm:"type:varchar(255);comment:邮箱"`
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
	uid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.UserID = strings.ReplaceAll(uid.String(), "-", "")

	if err := u.Encrypt(); err != nil {
		return err
	}

	return nil
}

func (u *User) BeforeUpdate(_ *gorm.DB) error {
	if err := u.Encrypt(); err != nil {
		return err
	}
	return nil
}

func (u *User) Encrypt() error {
	if u.Password == "" {
		return nil
	}

	var hash []byte
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
