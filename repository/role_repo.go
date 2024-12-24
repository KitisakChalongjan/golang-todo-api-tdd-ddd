package repository

import (
	"golang-todo-api-tdd-ddd/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	CreateRole(name string) (domain.Role, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (repo *RoleRepository) CreateRole(name string) (domain.Role, error) {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Error
	if err != nil {
		return domain.Role{}, err
	}

	role := domain.Role{}

	role.ID = uuid.NewString()
	role.Name = name

	err = tx.Create(&role).Error
	if err != nil {
		tx.Rollback()
		return domain.Role{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.Role{}, err
	}

	return role, nil
}
