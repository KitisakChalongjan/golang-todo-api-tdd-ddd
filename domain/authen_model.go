package domain

import (
	"time"
)

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshToken struct {
	ID           string    `json:"id" column:"id" gorm:"primaryKey"`
	UserID       string    `json:"user_id" column:"user_id" gorm:"not null"`
	RefreshToken string    `json:"refresh_token" column:"refresh_token" gorm:"not null"`
	IsRevoked    bool      `json:"is_revoked" column:"is_revoked" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" column:"created_at" gorm:"autoCreateTime:milli;not null"`
	UpdatedAt    time.Time `json:"updated_at" column:"updated_at" gorm:"autoUpdateTime:milli;not null"`
	DeviceInfo   string    `json:"device_info" column:"device_info" gorm:"not null"`
	IpAddress    string    `json:"ip_address" column:"ip_address" gorm:"not null"`
	User         User      `gorm:"constraint:OnDelete:CASCADE;"`
}

type UpdateRefreshTokenDTO struct {
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	IsRevoked    bool   `json:"is_revoked"`
	DeviceInfo   string `json:"device_info"`
	IpAddress    string `json:"ip_address"`
}
