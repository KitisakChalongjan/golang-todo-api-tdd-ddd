package service

import (
	"errors"
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthenService struct {
	userRepo   *repository.UserRepository
	authenRepo *repository.AuthenRepository
}

func NewAuthenService(userRepo *repository.UserRepository, authenRepo *repository.AuthenRepository) *AuthenService {
	return &AuthenService{userRepo: userRepo, authenRepo: authenRepo}
}

func (service *AuthenService) SignUpUser(user *domain.User, userDTO domain.SignUpUserDTO) error {

	bytes, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 14)
	if err != nil {
		return err
	}

	userDTO.Password = string(bytes)

	if err := service.authenRepo.CreateUser(user, userDTO); err != nil {
		return err
	}

	return nil
}

func (service *AuthenService) Login(accessToken *string, refreshToken *string, loginDTO domain.LoginDTO) error {

	userDTO := domain.GetUserDTO{}

	if err := service.authenRepo.GetUserByCredential(&userDTO, loginDTO); err != nil {
		return err
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return errors.New("JWT_SECRET environment variable not set")
	}

	jwtClaims := jwt.RegisteredClaims{
		Subject:   userDTO.ID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	signedAccessJWT, err := accessJWT.SignedString([]byte(secretKey))
	if err != nil {
		return fmt.Errorf("cannot creating access jwt. error : %s", err)
	}

	*accessToken = signedAccessJWT

	issuedAt := jwt.NewNumericDate(time.Now())
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30))

	refreshClaims := jwt.RegisteredClaims{
		Subject:   userDTO.ID,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedRefreshJWT, err := refreshJWT.SignedString([]byte(secretKey))
	if err != nil {
		return fmt.Errorf("cannot creating refresh jwt. error : %s", err)
	}

	*refreshToken = signedRefreshJWT

	updateRefreshTokenDTO := domain.UpdateRefreshTokenDTO{
		UserID:       userDTO.ID,
		RefreshToken: *refreshToken,
		IsRevoked:    false,
		DeviceInfo:   loginDTO.DeviceInfo,
		IpAddress:    loginDTO.IpAddress,
		IssuedAt:     issuedAt.Time,
		ExpiresAt:    expiresAt.Time,
	}

	if err := service.authenRepo.UpdateRefreshToken(updateRefreshTokenDTO); err != nil {
		return err
	}

	return nil
}

func (service *AuthenService) Logout(accessToken *jwt.Token) error {

	var getUserDTO domain.GetUserDTO

	userID, err := accessToken.Claims.GetSubject()
	if err != nil {
		return fmt.Errorf("cannot get userID from jwt. error : %s", err)
	}

	if err := service.userRepo.GetUserById(&getUserDTO, userID); err != nil {
		return err
	}

	issuedAt := jwt.NewNumericDate(time.Now())

	updateRefreshTokenDTO := domain.UpdateRefreshTokenDTO{
		UserID:       getUserDTO.ID,
		RefreshToken: "",
		IsRevoked:    true,
		DeviceInfo:   "",
		IpAddress:    "",
		IssuedAt:     issuedAt.Time,
		ExpiresAt:    issuedAt.Time,
	}

	if err := service.authenRepo.UpdateRefreshToken(updateRefreshTokenDTO); err != nil {
		return err
	}

	return nil
}
