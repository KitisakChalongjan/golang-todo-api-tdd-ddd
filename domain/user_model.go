package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID      `json:"id" column:"id" gorm:"primaryKey"`
	Name       string         `json:"name" column:"name" gorm:"not null"`
	Email      string         `json:"email" column:"email" gorm:"not null"`
	ProfileURL sql.NullString `json:"profile_url" column:"profile_url" gorm:""`
	CreatedAt  time.Time      `json:"created_at" column:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt  time.Time      `json:"updated_at" column:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt  sql.NullTime   `json:"deleted_at" column:"deleted_at" gorm:""`
}
