package model

type RoleResource struct {
	RoleID     int64 `json:"role_id" gorm:"not null;index:role_id_resource_id_index;comment:角色ID"`
	ResourceID int64 `json:"resource_id" gorm:"not null;index:role_id_resource_id_index;comment:资源ID"`
}

func (r *RoleResource) TableName() string {
	return "sys_role_resource"
}
