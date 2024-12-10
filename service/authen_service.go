package service

import (
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
)

type AuthenService struct {
	userRepo *repository.UserRepository
}

func NewAuthenService(userRepo *repository.UserRepository) *AuthenService {
	return &AuthenService{userRepo: userRepo}
}

func (handler *AuthenService) Login(userDTO *domain.GetUserDTO, logingDTO domain.LogingDTO) error {

	if err := handler.userRepo.GetUserByCredential(userDTO, logingDTO); err != nil {
		return err
	}

	return nil
}
