package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleAdmin = "admin" // 管理员角色 Code
	RoleUser  = "user"  // 普通用户角色 Code
)

type Role struct {
	ID        int64          `json:"-" gorm:"autoIncrement;primaryKey"`
	RoleID    string         `json:"role_id" gorm:"varchar(36);unique;not null;comment:角色ID"`
	Code      string         `json:"code" gorm:"varchar(100);uniqueIndex;not null;comment:角色标识"`
	Name      string         `json:"name" gorm:"varchar(255);not null;comment:角色名称"`
	Status    EnabledStatus  `json:"status" gorm:"not null;default:1;comment:状态,1:启用,2:禁用"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

func (r *Role) TableName() string {
	return "sys_role"
}

func (r *Role) BeforeCreate(_ *gorm.DB) error {
	uid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	r.RoleID = strings.ReplaceAll(uid.String(), "-", "")
	return nil
}

func (r *Role) Enabled() bool {
	return r.Status == Enabled
}
