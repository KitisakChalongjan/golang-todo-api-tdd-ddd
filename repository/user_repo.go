package repository

import (
	"golang-todo-api-tdd-ddd/domain"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetAllUsers(users *[]domain.User) *gorm.DB
	GetUserByID(user *domain.User, userID string) *gorm.DB
	CreateUser(user *domain.User) *gorm.DB
	UpdateUser(user *domain.User) *gorm.DB
	DeleteUser(userID string) *gorm.DB
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (userRepo *UserRepository) GetAllUsers(users *[]domain.User) *gorm.DB {

	result := userRepo.db.Find(&users)

	return result
}

func (userRepo *UserRepository) GetUserByID(user *domain.User, userID string) *gorm.DB {

	result := userRepo.db.Where("id = ?", userID).First(&user)

	return result
}

func (userRepo *UserRepository) CreateUser(user *domain.User) *gorm.DB {

	result := userRepo.db.Create(&user)

	return result
}

func (userRepo *UserRepository) UpdateUser(user *domain.User) *gorm.DB {

	result := userRepo.db.Save(&user)

	return result
}

func (userRepo *UserRepository) DeleteUser(userID string) *gorm.DB {

	result := userRepo.db.Where("id = ?", userID).Delete(&domain.User{})

	return result
}
