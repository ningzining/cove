package model

type RoleResource struct {
	RoleID     string `json:"role_id" gorm:"type:varchar(64);not null;index:role_id_resource_id_index;comment:角色ID"`
	ResourceID string `json:"resource_id" gorm:"type:varchar(64);not null;index:role_id_resource_id_index;comment:资源ID"`
}

func (r *RoleResource) TableName() string {
	return "sys_role_resource"
}
