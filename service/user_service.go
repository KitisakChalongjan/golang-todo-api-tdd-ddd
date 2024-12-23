package service

import (
	"fmt"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/valueobject"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (service *UserService) GetUserByID(userID string, accessToken *jwt.Token) (valueobject.GetUserVO, error) {

	tokenUserId, err := accessToken.Claims.GetSubject()
	if err != nil {
		return valueobject.GetUserVO{}, fmt.Errorf("fail to get userID from access token: %w", err)
	}

	getUserVO, err := service.userRepo.GetUserById(userID, tokenUserId)
	if err != nil {
		return valueobject.GetUserVO{}, err
	}

	return getUserVO, nil
}

func (service *UserService) UpdateUser(updateUserVO *valueobject.UpdateUserVO, accessToken *jwt.Token) (string, error) {

	tokenUserId, err := accessToken.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("fail to get userID from access token: %w", err)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(updateUserVO.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("fail to generate hash from passwrd: %w", err)
	}

	updateUserVO.Password = string(bytes)

	user, err := service.userRepo.UpdateUser(updateUserVO, tokenUserId)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (service *UserService) DeleteUser(userID string, accessToken *jwt.Token) (string, error) {

	tokenUserId, err := accessToken.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("fail to get userID from access token: %w", err)
	}

	user, err := service.userRepo.DeleteUser(userID, tokenUserId)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}
