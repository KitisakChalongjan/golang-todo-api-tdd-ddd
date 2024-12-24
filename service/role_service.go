package service

import (
	"fmt"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/valueobject"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (service *RoleService) CreateRole(roleVO valueobject.CreateRoleVO) (string, error) {

	role, err := service.roleRepo.CreateRole(roleVO.Name)
	if err != nil {
		return "", fmt.Errorf("create new role fail: %w", err)
	}

	return role.ID, nil
}
