package domain

type TransactionType struct {
	ID   string `json:"id" column:"id" gorm:"primaryKey"`
	Name string `json:"name" column:"name" gorm:"uniqueIndex;not null"`
}
