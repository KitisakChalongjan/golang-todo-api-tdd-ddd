package repository

import (
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthenRepository interface {
	CreateUser(user *domain.User, userDTO domain.SignUpUserDTO) error
	GetUserByCredential(loginDTO domain.LoginDTO) error
	UpdateRefreshToken(updateRefreshTokenDTO domain.UpdateRefreshTokenDTO) error
}

type AuthenRepository struct {
	db *gorm.DB
}

func NewAuthenRepository(db *gorm.DB) *AuthenRepository {
	return &AuthenRepository{db: db}
}

// create user
func (repo *AuthenRepository) CreateUser(user *domain.User, userDTO domain.SignUpUserDTO) error {

	var refreshToken domain.RefreshToken

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
	user.Role = userDTO.Role
	user.ProfileImgURL = userDTO.ProfileImgURL
	user.Username = userDTO.Username
	user.PasswordHash = userDTO.Password

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create user fail. error : %s", err.Error())
	}

	refreshToken.ID = uuid.New().String()
	refreshToken.UserID = user.ID
	refreshToken.RefreshToken = ""
	refreshToken.IsRevoked = true
	refreshToken.DeviceInfo = ""
	refreshToken.IpAddress = ""
	refreshToken.IssuedAt = time.Now()
	refreshToken.ExpiresAt = time.Now()

	if err := tx.Save(&refreshToken).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create refreshToken fail. error : %s", err.Error())
	}

	return tx.Commit().Error
}

func (repo *AuthenRepository) GetUserByCredential(userDTO *domain.GetUserDTO, loginDTO domain.LoginDTO) error {

	user := domain.User{}

	if err := repo.db.Where("username = ?", loginDTO.Username).First(&user).Error; err != nil {
		return fmt.Errorf("user not found. error : %s", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDTO.Password)); err != nil {
		return fmt.Errorf("incorrect password. error : %s", err.Error())
	}

	userDTO.ID = user.ID
	userDTO.Name = user.Name
	userDTO.Email = user.Email
	userDTO.Username = user.Username
	userDTO.ProfileImgURL = user.ProfileImgURL
	userDTO.CreatedAt = user.CreatedAt

	return nil
}

func (repo *AuthenRepository) UpdateRefreshToken(updateRefreshTokenDTO domain.UpdateRefreshTokenDTO) error {

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

	if err := tx.Where("user_id = ?", updateRefreshTokenDTO.UserID).First(refreshToken).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("no refresh token found. error : %s", err.Error())
	}

	refreshToken.RefreshToken = updateRefreshTokenDTO.RefreshToken
	refreshToken.IsRevoked = updateRefreshTokenDTO.IsRevoked
	refreshToken.DeviceInfo = updateRefreshTokenDTO.DeviceInfo
	refreshToken.IpAddress = updateRefreshTokenDTO.IpAddress
	refreshToken.IssuedAt = updateRefreshTokenDTO.IssuedAt
	refreshToken.ExpiresAt = updateRefreshTokenDTO.ExpiresAt

	if err := tx.Save(&refreshToken).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("update refreshToken fail. error : %s", err.Error())
	}

	return tx.Commit().Error
}
