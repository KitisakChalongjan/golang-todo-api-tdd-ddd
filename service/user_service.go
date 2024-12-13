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

func (service *UserService) GetAllUser(allUserDTO *[]domain.GetUserDTO) error {

	if err := service.userRepo.GetAllUsers(allUserDTO); err != nil {
		return err
	}

	return nil
}

func (service *UserService) GetUser(user *domain.GetUserDTO, userID string) error {

	if err := service.userRepo.GetUserById(user, userID); err != nil {
		return err
	}

	return nil
}

func (service *UserService) UpdateUser(user *domain.User, userDTO *domain.UpdateUserDTO) error {

	bytes, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userDTO.Password = string(bytes)

	err = service.userRepo.UpdateUser(user, userDTO)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) DeleteUser(userID string) error {

	err := service.userRepo.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}
