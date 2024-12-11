package service

import (
	"errors"
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthenService struct {
	userRepo *repository.UserRepository
}

func NewAuthenService(userRepo *repository.UserRepository) *AuthenService {
	return &AuthenService{userRepo: userRepo}
}

func (handler *AuthenService) Login(tokenString *string, logingDTO domain.LogingDTO) error {

	userDTO := domain.GetUserDTO{}

	if err := handler.userRepo.GetUserByCredential(&userDTO, logingDTO); err != nil {
		return err
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return errors.New("JWT_SECRET environment variable not set")
	}

	claims := jwt.MapClaims{
		"user_id": userDTO.ID,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return fmt.Errorf("cannot creating jwt token. error : %s", err)
	}

	*tokenString = signedString

	return nil
}
