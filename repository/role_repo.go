package repository

import (
	"golang-todo-api-tdd-ddd/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	CreateRole(name string) (*domain.Role, error)
	GetRoleIDsByRoleNames(roleNames []string) (*[]domain.Role, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (repo *RoleRepository) CreateRole(name string) (*domain.Role, error) {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Error
	if err != nil {
		return nil, err
	}

	newRole := domain.Role{}

	newRole.ID = uuid.NewString()
	newRole.Name = name

	err = tx.Create(&newRole).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &newRole, nil
}
