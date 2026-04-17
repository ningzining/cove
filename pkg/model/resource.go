package model

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActionList []string

func (a ActionList) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ActionList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

type Resource struct {
	ID         int64          `json:"-" gorm:"autoIncrement;primaryKey"`
	ResourceID string         `json:"resource_id" gorm:"varchar(36);unique;not null;comment:资源ID"`
	Code       string         `json:"code" gorm:"varchar(100);unique;not null;comment:资源标识"`
	Name       string         `json:"name" gorm:"varchar(255);not null;comment:资源名称"`
	Actions    ActionList     `json:"actions" gorm:"type:json;comment:支持的操作列表"`
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
