package domain

import (
	"time"
)

type Transaction struct {
	ID                string          `json:"id" column:"id" gorm:"primaryKey"`
	UserID            string          `json:"user_id" column:"user_id" gorm:"not null"`
	TransactionTypeID string          `json:"transaction_type_id" column:"transaction_type_id" gorm:"not null"`
	Amount            float32         `json:"amount" column:"amount" gorm:"not null"`
	Description       string          `json:"description" column:"description" gorm:"not null"`
	CreatedAt         time.Time       `json:"created_at" column:"created_at" gorm:"autoCreateTime:milli;not null"`
	UpdatedAt         time.Time       `json:"updated_at" column:"updated_at" gorm:"autoUpdateTime:milli;not null"`
	DeletedAt         *time.Time      `json:"deleted_at" column:"deleted_at" gorm:""`
	User              User            `gorm:"constraint:OnDelete:CASCADE"`
	TransactionType   TransactionType `gorm:"constraint:OnDelete:CASCADE"`
}
