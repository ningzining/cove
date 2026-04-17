package model

type UserRole struct {
	UserID int64 `json:"user_id" gorm:"not null;index:user_id_role_id_index;comment:用户ID"`
	RoleID int64 `json:"role_id" gorm:"not null;index:user_id_role_id_index;comment:角色ID"`
}

func (u *UserRole) TableName() string {
	return "sys_user_role"
}
