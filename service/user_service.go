package service

import (
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (userService *UserService) GetAllUser(allUserDTO *[]domain.GetUserDTO) error {

	if err := userService.userRepo.GetAllUsers(allUserDTO); err != nil {
		return err
	}

	return nil
}

func (userService *UserService) GetUser(user *domain.GetUserDTO, userID string) error {

	if err := userService.userRepo.GetUserById(user, userID); err != nil {
		return err
	}

	return nil
}

func (userService *UserService) CreateUser(user *domain.User, userDTO domain.CreateUserDTO) error {

	bytes, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userDTO.Password = string(bytes)

	if err := userService.userRepo.CreateUser(user, userDTO); err != nil {
		return err
	}

	return nil
}

func (userService *UserService) UpdateUser(user *domain.User, userDTO *domain.UpdateUserDTO) error {

	bytes, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userDTO.Password = string(bytes)

	err = userService.userRepo.UpdateUser(user, userDTO)
	if err != nil {
		return err
	}

	return nil
}

func (userService *UserService) DeleteUser(userID string) error {

	err := userService.userRepo.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}
