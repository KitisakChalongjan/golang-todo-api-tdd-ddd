package domain

import (
	"time"
)

type User struct {
	ID            string     `json:"id" column:"id" gorm:"primaryKey"`
	Name          string     `json:"name" column:"name" gorm:"not null"`
	Email         string     `json:"email" column:"email" gorm:"uniqueIndex;not null"`
	ProfileImgURL *string    `json:"profile_img_url" column:"profile_img_url" gorm:""`
	Username      string     `json:"username" column:"username" gorm:"uniqueIndex;not null"`
	PasswordHash  string     `json:"password_hash" column:"password_hash" gorm:"not null"`
	Roles         []Role     `json:"roles" column:"roles" gorm:"many2many:users_roles;not null"`
	CreatedAt     time.Time  `json:"created_at" column:"created_at" gorm:"autoCreateTime:milli;not null"`
	UpdatedAt     time.Time  `json:"updated_at" column:"updated_at" gorm:"autoUpdateTime:milli;not null"`
	DeletedAt     *time.Time `json:"deleted_at" column:"deleted_at" gorm:""`
}
