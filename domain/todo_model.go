package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID    `json:"id" column:"id" gorm:"primaryKey"`
	Title       string       `json:"title" column:"title" gorm:"not null"`
	Description string       `json:"description" column:"description" gorm:"not null"`
	IsCompleted bool         `json:"is_complete" column:"is_completed" gorm:"not null"`
	Priority    string       `json:"priority" column:"priority" gorm:"not null"`
	Due         sql.NullTime `json:"due" column:"due" gorm:""`
	CreatedAt   time.Time    `json:"created_at" column:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt   time.Time    `json:"updated_at" column:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt   sql.NullTime `json:"deleted_at" column:"deleted_at" gorm:""`
	UserID      uuid.UUID    `json:"user_id" column:"user_id" gorm:"not null"`
	User        User         `gorm:"constraint:OnDelete:CASCADE;"`
}
