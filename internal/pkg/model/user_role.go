package model

type UserRole struct {
	UserID string `json:"user_id" gorm:"type:varchar(64);not null;index:user_id_role_id_index;comment:用户ID"`
	RoleID string `json:"role_id" gorm:"type:varchar(64);not null;index:user_id_role_id_index;comment:角色ID"`
}

func (u *UserRole) TableName() string {
	return "sys_user_role"
}
