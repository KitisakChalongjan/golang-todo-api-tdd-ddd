package repository

import (
	"golang-todo-api-tdd-ddd/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetAllUsers(users *[]domain.User) error
	GetUserById(user *domain.User, userID string) error
	CreateUser(user *domain.User) error
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
		return err
	}

	*allUserDTO = make([]domain.GetUserDTO, len(allUser))

	for i, user := range allUser {
		(*allUserDTO)[i] = domain.GetUserDTO{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email,
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
	userDTO.ProfileImgURL = user.ProfileImgURL
	userDTO.Username = user.Username
	userDTO.CreatedAt = user.CreatedAt

	return nil
}

// create user
func (repo *UserRepository) CreateUser(user *domain.User, userDTO domain.CreateUserDTO) error {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	user.ID = uuid.New().String()
	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.ProfileImgURL = userDTO.ProfileImgURL
	user.Username = userDTO.Username
	user.PasswordHash = userDTO.Password

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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
		return err
	}

	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.ProfileImgURL = userDTO.ProfileImgURL
	user.PasswordHash = userDTO.Password

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
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
		return err
	}

	return tx.Commit().Error
}
