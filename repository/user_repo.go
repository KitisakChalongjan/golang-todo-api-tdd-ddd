package repository

import (
	"fmt"
	"golang-todo-api-tdd-ddd/domain"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetAllUsers(users *[]domain.User) error
	GetUserById(user *domain.User, userID string) error
	SignInUser(user *domain.User) error
	UpdateUser(userID string, userDTO domain.UpdateUserDTO) error
	DeleteUser(userID string) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) GetAllUsers(allUserDTO *[]domain.GetUserDTO) error {

	allUser := []domain.User{}

	if err := repo.db.Find(&allUser).Error; err != nil {
		return fmt.Errorf("no users found. error : %s", err.Error())
	}

	*allUserDTO = make([]domain.GetUserDTO, len(allUser))

	for i, user := range allUser {
		(*allUserDTO)[i] = domain.GetUserDTO{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email,
			Role:          user.Role,
			ProfileImgURL: user.ProfileImgURL,
			Username:      user.Username,
			CreatedAt:     user.CreatedAt,
		}
	}

	return nil
}

func (repo *UserRepository) GetUserById(userDTO *domain.GetUserDTO, userID string) error {

	user := domain.User{}

	if err := repo.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	userDTO.ID = user.ID
	userDTO.Name = user.Name
	userDTO.Email = user.Email
	userDTO.Role = user.Role
	userDTO.ProfileImgURL = user.ProfileImgURL
	userDTO.Username = user.Username
	userDTO.CreatedAt = user.CreatedAt

	return nil
}

// update user
func (repo *UserRepository) UpdateUser(user *domain.User, userDTO *domain.UpdateUserDTO) error {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("id = ?", userDTO.ID).First(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("no user found. error : %s", err.Error())
	}

	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.Role = userDTO.Role
	user.ProfileImgURL = userDTO.ProfileImgURL
	user.PasswordHash = userDTO.Password

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("update user fail. error : %s", err.Error())
	}

	return tx.Commit().Error
}

func (repo *UserRepository) DeleteUser(userID string) error {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("id = ?", userID).Delete(&domain.User{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("delete user fail. error : %s", err.Error())
	}

	return tx.Commit().Error
}
