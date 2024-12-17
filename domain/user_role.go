package domain

type UsersRoles struct {
	UserID string `json:"user_id" column:"user_id" gorm:"primaryKey"`
	RoleID string `json:"role_id" column:"role_id" gorm:"primaryKey"`
}
