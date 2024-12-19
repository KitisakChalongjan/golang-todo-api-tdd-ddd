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

// func (service *UserService) GetAllUser(allUserDTO *[]valueobject.GetUserVO) error {

// 	if err := service.userRepo.GetAllUsers(allUserDTO); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (service *UserService) GetUserByID(userID string, accessToken *jwt.Token) (valueobject.GetUserVO, error) {

	tokenUserId, err := accessToken.Claims.GetSubject()
	if err != nil {
		return valueobject.GetUserVO{}, fmt.Errorf("fail to get userID from access token: %w", err)
	}

	if tokenUserId != userID {
		return valueobject.GetUserVO{}, fmt.Errorf("you are not authorized to access this data: token userID and param userID not match")
	}

	getUserVO, err := service.userRepo.GetUserById(userID)
	if err != nil {
		return valueobject.GetUserVO{}, err
	}

	return getUserVO, nil
}

func (service *UserService) UpdateUser(updateUserDTO *valueobject.UpdateUserVO, accessToken *jwt.Token) (string, error) {

	tokenUserId, err := accessToken.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("fail to get userID from access token: %w", err)
	}

	if tokenUserId != updateUserDTO.ID {
		return "", fmt.Errorf("you are not authorized to access this data: token userID and param userID not match")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(updateUserDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("fail to generate hash from passwrd: %w", err)
	}

	updateUserDTO.Password = string(bytes)

	userID, err := service.userRepo.UpdateUser(updateUserDTO)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (service *UserService) DeleteUser(userID string) error {

	err := service.userRepo.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}
