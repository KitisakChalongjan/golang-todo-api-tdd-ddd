package domain

import (
	"time"
)

type Todo struct {
	ID          string     `json:"id" column:"id" gorm:"primaryKey"`
	UserID      string     `json:"user_id" column:"user_id" gorm:"not null"`
	Title       string     `json:"title" column:"title" gorm:"not null"`
	Description string     `json:"description" column:"description" gorm:"not null"`
	IsCompleted bool       `json:"is_completed" column:"is_completed" gorm:"not null"`
	Priority    string     `json:"priority" column:"priority" gorm:"not null"`
	Due         *time.Time `json:"due" column:"due" gorm:""`
	CreatedAt   time.Time  `json:"created_at" column:"created_at" gorm:"autoCreateTime:milli;not null"`
	UpdatedAt   time.Time  `json:"updated_at" column:"updated_at" gorm:"autoUpdateTime:milli;not null"`
	DeletedAt   *time.Time `json:"deleted_at" column:"deleted_at" gorm:""`
	User        User       `gorm:"constraint:OnDelete:CASCADE;"`
}

type GetTodoDTO struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	Priority    string     `json:"priority"`
	Due         *time.Time `json:"due"`
	UserID      string     `json:"user_id"`
}

type CreateTodoDTO struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	Due         *time.Time `json:"due"`
	UserID      string     `json:"user_id"`
}

type UpdateTodoDTO struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	Priority    string     `json:"priority"`
	Due         *time.Time `json:"due"`
}
