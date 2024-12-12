package service

import (
	"errors"
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthenService struct {
	userRepo   *repository.UserRepository
	authenRepo *repository.AuthenRepository
}

func NewAuthenService(userRepo *repository.UserRepository, authenRepo *repository.AuthenRepository) *AuthenService {
	return &AuthenService{userRepo: userRepo, authenRepo: authenRepo}
}

func (handler *AuthenService) Login(accessToken *string, refreshToken *string, loginDTO domain.LoginDTO, c echo.Context) error {

	userDTO := domain.GetUserDTO{}

	if err := handler.authenRepo.GetUserByCredential(&userDTO, loginDTO); err != nil {
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

	refreshClaims := jwt.RegisteredClaims{
		Subject:   userDTO.ID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
	}

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedRefreshJWT, err := refreshJWT.SignedString([]byte(secretKey))
	if err != nil {
		return fmt.Errorf("cannot creating refresh jwt. error : %s", err)
	}

	*refreshToken = signedRefreshJWT

	// updateRefreshTokenDTO := domain.UpdateRefreshTokenDTO{
	// 	UserID:       userDTO.ID,
	// 	RefreshToken: *refreshToken,
	// 	IsRevoked:    false,
	// 	DeviceInfo:   "",
	// 	IpAddress:    "",
	// }

	handler.authenRepo.UpdateRefreshToken(userDTO.ID, *refreshToken)

	return nil
}
