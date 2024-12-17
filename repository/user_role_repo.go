package repository

import (
	"golang-todo-api-tdd-ddd/domain"

	"gorm.io/gorm"
)

type IUserRoleRepository interface {
	CreateUserRoles(userID string, roleIDs []string) (*[]domain.UsersRoles, error)
}

type UserRoleRepository struct {
	DB *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{DB: db}
}

func (repo *UserRoleRepository) CreateUserRoles(userID string, roleIDs []string) ([]domain.UsersRoles, error) {
	tx := repo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	newUserRoles := []domain.UsersRoles{}

	for _, roleID := range roleIDs {
		newUserRole := domain.UsersRoles{
			UserID: userID,
			RoleID: roleID,
		}
		newUserRoles = append(newUserRoles, newUserRole)
	}

	if err := tx.Create(&newUserRoles).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return newUserRoles, nil
}
