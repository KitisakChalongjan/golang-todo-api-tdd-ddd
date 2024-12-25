package repository

import (
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/valueobject"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByID(userID string, tokenUserID string) (valueobject.GetUserVO, error)
	GetUserByCredential(signInVO valueobject.SignInVO) (valueobject.GetUserVO, error)
	CreateUser(signUpVO valueobject.SignUpVO) (domain.User, error)
	UpdateUser(updateUserVO valueobject.UpdateUserVO, tokenUserID string) (domain.User, error)
	DeleteUser(userID string, tokenUserId string) (domain.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) GetUserByID(userID string, tokenUserID string) (valueobject.GetUserVO, error) {

	user := domain.User{}
	userVO := valueobject.GetUserVO{}

	if err := repo.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return valueobject.GetUserVO{}, fmt.Errorf("fail to get user by id.: %w", err)
	}

	if tokenUserID != user.ID {
		return valueobject.GetUserVO{}, fmt.Errorf("you cannot access this data: userID not match")
	}

	roleNames, err := GetRoleNamesByUserID(repo, user.ID)
	if err != nil {
		return valueobject.GetUserVO{}, err
	}

	userVO.ID = user.ID
	userVO.Name = user.Name
	userVO.Email = user.Email
	userVO.ProfileImgURL = user.ProfileImgURL
	userVO.Username = user.Username
	userVO.Roles = roleNames
	userVO.CreatedAt = user.CreatedAt

	return userVO, nil
}

func (repo *UserRepository) GetUserByCredential(signInVO valueobject.SignInVO) (valueobject.GetUserVO, error) {

	user := domain.User{}
	getUserVO := valueobject.GetUserVO{}

	if err := repo.db.Where("username = ?", signInVO.Username).First(&user).Error; err != nil {
		return valueobject.GetUserVO{}, fmt.Errorf("cannot find user from username %s. error:  %w", signInVO.Username, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(signInVO.Password)); err != nil {
		return valueobject.GetUserVO{}, fmt.Errorf("invalid password for username '%s'. error: %w", signInVO.Username, err)
	}

	roleNames, err := GetRoleNamesByUserID(repo, user.ID)
	if err != nil {
		return valueobject.GetUserVO{}, err
	}

	getUserVO.ID = user.ID
	getUserVO.Name = user.Name
	getUserVO.Email = user.Email
	getUserVO.Username = user.Username
	getUserVO.ProfileImgURL = user.ProfileImgURL
	getUserVO.Roles = roleNames
	getUserVO.CreatedAt = user.CreatedAt

	return getUserVO, nil
}

func (repo *UserRepository) CreateUser(signUpVO valueobject.SignUpVO) (domain.User, error) {

	tx := repo.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return domain.User{}, fmt.Errorf("transaction fail: %s", err.Error())
	}

	user := domain.User{}
	user.ID = uuid.New().String()
	user.Name = signUpVO.Name
	user.Email = signUpVO.Email
	user.ProfileImgURL = signUpVO.ProfileImgURL
	user.Username = signUpVO.Username
	user.PasswordHash = signUpVO.Password

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	roles := []domain.Role{}
	err := repo.db.Where("name IN (?)", signUpVO.Roles).Find(&roles).Error
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to find roles: %w", err)
	}

	userRoles := make([]domain.UsersRoles, len(roles))
	for i, role := range roles {
		userRoles[i] = domain.UsersRoles{
			UserID: user.ID,
			RoleID: role.ID,
		}
	}

	if err := tx.Create(&userRoles).Error; err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("failed to create user roles: %w", err)
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return user, nil
}

func (repo *UserRepository) UpdateUser(updateUserVO valueobject.UpdateUserVO, tokenUserID string) (domain.User, error) {

	user := domain.User{}

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return domain.User{}, fmt.Errorf("transaction fail: %s", err.Error())
	}

	err := tx.Where("id = ?", updateUserVO.ID).First(&user).Error
	if err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("no user found: %s", err.Error())
	}

	if tokenUserID != user.ID {
		return domain.User{}, fmt.Errorf("you cannot access this data: userID not match")
	}

	err = tx.Where("user_id = ?", updateUserVO.ID).Delete(&domain.UsersRoles{}).Error
	if err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("fail to update user roles: %s", err.Error())
	}

	roles := []domain.Role{}

	err = tx.Where("name IN (?)", updateUserVO.Roles).Find(&roles).Error
	if err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("no roles found: %s", err.Error())
	}

	newUserRoles := make([]domain.UsersRoles, len(roles))
	for i, role := range roles {
		newUserRoles[i] = domain.UsersRoles{
			UserID: updateUserVO.ID,
			RoleID: role.ID,
		}
	}

	if err := tx.Create(&newUserRoles).Error; err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("failed to create user roles: %w", err)
	}

	user.Name = updateUserVO.Name
	user.Email = updateUserVO.Email
	user.ProfileImgURL = updateUserVO.ProfileImgURL
	user.PasswordHash = updateUserVO.Password

	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("update user fail: %s", err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.User{}, fmt.Errorf("update user commit fail: %s", err.Error())
	}

	return user, nil
}

func (repo *UserRepository) DeleteUser(userID string, tokenUserId string) (domain.User, error) {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := domain.User{}

	if err := tx.Error; err != nil {
		return domain.User{}, fmt.Errorf("transaction fail: %s", err.Error())
	}

	err := tx.Where("id = ?", userID).First(&user).Error
	if err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("no user found: %s", err.Error())
	}

	if tokenUserId != user.ID {
		return domain.User{}, fmt.Errorf("you cannot access this data: userID not match")
	}

	err = tx.Where("user_id = ?", userID).Delete(&domain.UsersRoles{}).Error
	if err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("fail to delete user roles: %s", err.Error())
	}

	err = tx.Where("id = ?", userID).Delete(&domain.User{}).Error
	if err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("delete user fail: %s", err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.User{}, fmt.Errorf("delete user commit fail: %s", err.Error())
	}

	return user, nil
}

func GetRoleNamesByUserID(repo *UserRepository, userID string) ([]string, error) {

	usersRoles := []domain.UsersRoles{}

	err := repo.db.Where("user_id = ?", userID).Find(&usersRoles).Error
	if err != nil {
		return []string{}, fmt.Errorf("fail to get user roles by userID. error : %w", err)
	}

	roleIDs := make([]string, len(usersRoles))
	for i, userRole := range usersRoles {
		roleIDs[i] = userRole.RoleID
	}

	roles := []domain.Role{}
	err = repo.db.Where("id IN (?)", roleIDs).Find(&roles).Error
	if err != nil {
		return []string{}, fmt.Errorf("fail to get roles by roleIDs. error : %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return roleNames, nil
}
