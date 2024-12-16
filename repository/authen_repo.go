package repository

import (
	"fmt"
	"golang-todo-api-tdd-ddd/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthenRepository interface {
	CreateUser(signupDTO domain.SignUpDTO) (domain.User, error)
	GetUserByCredential(loginDTO domain.SignInDTO) error
	UpdateRefreshToken(updateRefreshTokenDTO domain.UpdateRefreshTokenDTO) (domain.GetUserDTO, error)
}

type AuthenRepository struct {
	db *gorm.DB
}

func NewAuthenRepository(db *gorm.DB) *AuthenRepository {
	return &AuthenRepository{db: db}
}

func (repo *AuthenRepository) CreateUser(signupDTO domain.SignUpDTO) (domain.User, error) {

	user := domain.User{}

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return domain.User{}, err
	}

	user.ID = uuid.New().String()
	user.Name = signupDTO.Name
	user.Email = signupDTO.Email
	user.Role = signupDTO.Role
	user.ProfileImgURL = signupDTO.ProfileImgURL
	user.Username = signupDTO.Username
	user.PasswordHash = signupDTO.Password

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return domain.User{}, err
	}

	err := tx.Commit().Error
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (repo *AuthenRepository) GetUserByCredential(loginDTO domain.SignInDTO) (domain.GetUserDTO, error) {

	user := domain.User{}
	getUserDTO := domain.GetUserDTO{}

	if err := repo.db.Where("username = ?", loginDTO.Username).First(&user).Error; err != nil {
		return domain.GetUserDTO{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDTO.Password)); err != nil {
		return domain.GetUserDTO{}, err
	}

	getUserDTO.ID = user.ID
	getUserDTO.Name = user.Name
	getUserDTO.Email = user.Email
	getUserDTO.Username = user.Username
	getUserDTO.ProfileImgURL = user.ProfileImgURL
	getUserDTO.CreatedAt = user.CreatedAt

	return getUserDTO, nil
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
