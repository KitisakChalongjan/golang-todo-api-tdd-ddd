package repository

import (
	"golang-todo-api-tdd-ddd/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetAllUsers() (*[]domain.User, error)
	GetUserByID(userID uuid.UUID) (*domain.User, error)
	CreateUser(user *domain.User) (*domain.User, error)
	UpdateUser(user *domain.User) (*domain.User, error)
	DeleteUser(userID uuid.UUID) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (userRepo *UserRepository) GetAllUsers() (*[]domain.User, error) {
	var users []domain.User

	result := userRepo.db.Find(&users)

	return &users, result.Error
}

func (userRepo *UserRepository) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	var user domain.User

	result := userRepo.db.First(&user, userID)

	return &user, result.Error
}

func (userRepo *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	result := userRepo.db.Create(user)

	return user, result.Error
}

func (userRepo *UserRepository) UpdateUser(user *domain.User) (*domain.User, error) {
	result := userRepo.db.Save(user)

	return user, result.Error
}

func (userRepo *UserRepository) DeleteUser(userID uuid.UUID) error {
	result := userRepo.db.Delete(&domain.User{}, userID)

	return result.Error
}
