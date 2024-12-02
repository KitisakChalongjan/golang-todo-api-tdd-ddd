package repository

import (
	"golang-todo-api-tdd-ddd/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetAllUsers(users *[]domain.User) error
	GetUser(user *domain.User, userID string) error
	CreateUser(user *domain.User) error
	UpdateUser(userID string, userDTO *domain.UpdateUserDTO) error
	DeleteUser(userID string) error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (userRepo *UserRepository) GetAllUsers(users *[]domain.User) error {

	// get all users
	if err := userRepo.DB.Find(&users).Error; err != nil {
		return err
	}

	return nil
}

func (userRepo *UserRepository) GetUser(user *domain.User, userID string) error {

	// get user by id
	if err := userRepo.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	return nil
}

// create user
func (userRepo *UserRepository) CreateUser(user *domain.User, userDTO *domain.CreateUserDTO) error {

	// begin transaction
	tx := userRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// check transaction error
	if err := tx.Error; err != nil {
		return err
	}

	user.ID = uuid.New().String()
	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.ProfileImgURL = userDTO.ProfileImgURL

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// commit transaction when no error
	return tx.Commit().Error
}

// update user
func (userRepo *UserRepository) UpdateUser(user *domain.User, userDTO *domain.UpdateUserDTO) error {

	// begin transaction
	tx := userRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// check transaction error
	if err := tx.Error; err != nil {
		return err
	}

	// get user by id
	if err := tx.Where("id = ?", userDTO.ID).First(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// assign userDTO to user
	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.ProfileImgURL = userDTO.ProfileImgURL

	// save user
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// commit transaction when no error
	return tx.Commit().Error
}

func (userRepo *UserRepository) DeleteUser(userID string) error {

	// begin transaction
	tx := userRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// check transaction error
	if err := tx.Error; err != nil {
		return err
	}

	// delete user by id
	if err := tx.Where("id = ?", userID).Delete(&domain.User{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
