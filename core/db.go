package core

import (
	"golang-todo-api-tdd-ddd/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgresDB() (*gorm.DB, error) {

	connectionString := "user=postgres password=Dewsmaller1* dbname=postgres host=localhost port=5432 sslmode=disable"

	dialector := postgres.Open(connectionString)

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("database connected.")

	err = db.AutoMigrate(&domain.Todo{}, &domain.User{}, &domain.RefreshToken{})
	if err != nil {
		return nil, err
	}

	log.Println("database migrate successful.")

	return db, nil
}
