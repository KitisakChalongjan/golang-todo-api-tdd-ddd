package core

import (
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgresDB() (*gorm.DB, error) {

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)

	dialector := postgres.Open(connectionString)

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Printf("database connected(DB_HOST=%s)", DB_HOST)

	err = db.AutoMigrate(
		&domain.Todo{},
		&domain.User{},
		&domain.Role{},
		&domain.UsersRoles{},
		&domain.Transaction{},
		&domain.TransactionType{},
	)
	if err != nil {
		return nil, err
	}

	log.Println("database migrate successful.")

	return db, nil
}
