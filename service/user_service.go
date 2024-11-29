package service

import (
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (userService *UserService) GetAllUser(users *[]domain.User) error {

	result := userService.userRepo.GetAllUsers(users)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (userService *UserService) GetUserByID(user *domain.User, userID string) error {

	result := userService.userRepo.GetUserByID(user, userID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (userService *UserService) CreateUser(user *domain.User, userDTO domain.CreateUserDTO) error {

	user.ID = uuid.New().String()
	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.ProfileImgURL = userDTO.ProfileImgURL

	result := userService.userRepo.CreateUser(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (userService *UserService) UpdateUser(user *domain.User, userDTO domain.UpdateUserDTO) error {

	result := userService.userRepo.GetUserByID(user, userDTO.ID)
	if result.Error != nil {
		return result.Error
	}

	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.ProfileImgURL = userDTO.ProfileImgURL

	result = userService.userRepo.UpdateUser(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (userService *UserService) DeleteUser(userID string) error {

	result := userService.userRepo.DeleteUser(userID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
