package repository

import (
	"errors"
	"golang-todo-api-tdd-ddd/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthenRepository interface {
	GetUserByCredential(loginDTO domain.LoginDTO) error
	UpdateRefreshToken(userID string, refreshToken string) error
}

type AuthenRepository struct {
	db *gorm.DB
}

func NewAuthenRepository(db *gorm.DB) *AuthenRepository {
	return &AuthenRepository{db: db}
}

func (repo *AuthenRepository) GetUserByCredential(userDTO *domain.GetUserDTO, loginDTO domain.LoginDTO) error {

	user := domain.User{}

	if err := repo.db.Where("username = ?", loginDTO.Username).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDTO.Password)); err != nil {
		return errors.New("incorrect password")
	}

	userDTO.ID = user.ID
	userDTO.Name = user.Name
	userDTO.Email = user.Email
	userDTO.Username = user.Username
	userDTO.ProfileImgURL = user.ProfileImgURL
	userDTO.CreatedAt = user.CreatedAt

	return nil
}

func (repo *AuthenRepository) UpdateRefreshToken(userID string, token string) error {

	refreshToken := &domain.RefreshToken{}

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("user_id = ?", userID).First(refreshToken).Error; err != nil {
		tx.Rollback()
		return err
	}

	refreshToken.RefreshToken = token
	refreshToken.IsRevoked = false

	if err := tx.Save(&refreshToken).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
