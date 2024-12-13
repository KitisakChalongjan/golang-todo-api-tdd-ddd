package domain

import (
	"time"
)

type SignUpUserDTO struct {
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Role          string  `json:"role"`
	ProfileImgURL *string `json:"profile_img_url"`
	Username      string  `json:"username"`
	Password      string  `json:"password"`
}

type LoginDTO struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceInfo string `json:"device_info"`
	IpAddress  string `json:"ip_address"`
}

type LogoutDTO struct {
	UserID string `json:"user_id"`
}

type RefreshToken struct {
	ID           string    `json:"id" column:"id" gorm:"primaryKey"`
	UserID       string    `json:"user_id" column:"user_id" gorm:"not null"`
	RefreshToken string    `json:"refresh_token" column:"refresh_token" gorm:"not null"`
	IsRevoked    bool      `json:"is_revoked" column:"is_revoked" gorm:"not null"`
	DeviceInfo   string    `json:"device_info" column:"device_info" gorm:"not null"`
	IpAddress    string    `json:"ip_address" column:"ip_address" gorm:""`
	CreatedAt    time.Time `json:"created_at" column:"created_at" gorm:"autoCreateTime:milli;not null"`
	IssuedAt     time.Time `json:"issued_at" column:"issued_at" gorm:"not null"`
	ExpiresAt    time.Time `json:"expires_at" column:"expires_at" gorm:"not null"`
	User         User      `gorm:"constraint:OnDelete:CASCADE;"`
}

type UpdateRefreshTokenDTO struct {
	UserID       string    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	IsRevoked    bool      `json:"is_revoked"`
	DeviceInfo   string    `json:"device_info"`
	IpAddress    string    `json:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at"`
	IssuedAt     time.Time `json:"issued_at"`
}

type UpdateAccessTokenDTO struct {
	UserID      string    `json:"user_id"`
	AccessToken string    `json:"refresh_token"`
	IsRevoked   bool      `json:"is_revoked"`
	DeviceInfo  string    `json:"device_info"`
	IpAddress   string    `json:"ip_address"`
	ExpiresAt   time.Time `json:"expires_at"`
	IssuedAt    time.Time `json:"issued_at"`
}

type AccessTokenClaims struct {
	UserID    string `json:"user_id"`
	IssuedAt  string `json:"iat"`
	ExpiresAt string `json:"exp"`
}
