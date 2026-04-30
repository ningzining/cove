package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Action string

const (
	ReadAction Action = "read" // 只读操作
	FullAction Action = "full" // 全部操作
)

type Resource struct {
	ID         int64          `json:"-" gorm:"autoIncrement;primaryKey"`
	ResourceID string         `json:"id" gorm:"type:varchar(64);unique;not null;comment:资源ID"`
	Code       string         `json:"code" gorm:"type:varchar(64);not null;comment:资源标识"`
	Name       string         `json:"name" gorm:"type:varchar(255);not null;comment:资源名称"`
	Action     Action         `json:"action" gorm:"type:varchar(64);comment:支持的操作"`
	CreatedAt  time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

func (r *Resource) TableName() string {
	return "sys_resource"
}

func (r *Resource) BeforeCreate(_ *gorm.DB) error {
	uid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	r.ResourceID = strings.ReplaceAll(uid.String(), "-", "")
	return nil
}
