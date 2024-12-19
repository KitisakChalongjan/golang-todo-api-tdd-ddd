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
	// GetAllUsers(users *[]domain.User) error
	GetUserById(userID string) (valueobject.GetUserVO, error)
	GetUserByCredential(loginDTO valueobject.SignInVO) error
	CreateUser(signupDTO valueobject.SignUpVO) (domain.User, error)
	SignInUser(user *domain.User) error
	UpdateUser(userID string, userDTO valueobject.UpdateUserVO) error
	DeleteUser(userID string) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// func (repo *UserRepository) GetAllUsers(allUserDTO *[]valueobject.GetUserVO) error {

// 	allUser := []domain.User{}

// 	if err := repo.db.Find(&allUser).Error; err != nil {
// 		return fmt.Errorf("no users found. error : %s", err.Error())
// 	}

// 	*allUserDTO = make([]valueobject.GetUserVO, len(allUser))

// 	for i, user := range allUser {
// 		(*allUserDTO)[i] = valueobject.GetUserVO{
// 			ID:            user.ID,
// 			Name:          user.Name,
// 			Email:         user.Email,
// 			ProfileImgURL: user.ProfileImgURL,
// 			Username:      user.Username,
// 			CreatedAt:     user.CreatedAt,
// 		}
// 	}

// 	return nil
// }

func (repo *UserRepository) GetUserById(userID string) (valueobject.GetUserVO, error) {

	user := domain.User{}
	userVO := valueobject.GetUserVO{}

	if err := repo.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return valueobject.GetUserVO{}, fmt.Errorf("fail to get user by id. error: %w", err)
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

func (repo *UserRepository) CreateUser(signupDTO valueobject.SignUpVO) (domain.User, error) {

	tx := repo.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return domain.User{}, err
	}

	newUser := domain.User{}
	newUser.ID = uuid.New().String()
	newUser.Name = signupDTO.Name
	newUser.Email = signupDTO.Email
	newUser.ProfileImgURL = signupDTO.ProfileImgURL
	newUser.Username = signupDTO.Username
	newUser.PasswordHash = signupDTO.Password

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	roles := []domain.Role{}
	err := repo.db.Where("name IN (?)", signupDTO.Roles).Find(&roles).Error
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to find roles: %w", err)
	}

	newUserRoles := make([]domain.UsersRoles, len(roles))
	for i, role := range roles {
		newUserRoles[i] = domain.UsersRoles{
			UserID: newUser.ID,
			RoleID: role.ID,
		}
	}

	if err := tx.Create(&newUserRoles).Error; err != nil {
		tx.Rollback()
		return domain.User{}, fmt.Errorf("failed to create user roles: %w", err)
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newUser, nil
}

func (repo *UserRepository) UpdateUser(updateUserDTO *valueobject.UpdateUserVO) (string, error) {

	user := domain.User{}

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return "", err
	}

	if err := tx.Where("id = ?", updateUserDTO.ID).First(&user).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("no user found: %s", err.Error())
	}

	roles := []domain.Role{}

	err := tx.Where("name IN (?)", updateUserDTO.Roles).Find(&roles).Error
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("no roles found: %s", err.Error())
	}

	user.Name = updateUserDTO.Name
	user.Email = updateUserDTO.Email
	user.Roles = roles
	user.ProfileImgURL = updateUserDTO.ProfileImgURL
	user.PasswordHash = updateUserDTO.Password

	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("update user fail: %s", err.Error())
	}

	return user.ID, tx.Commit().Error
}

func (repo *UserRepository) DeleteUser(userID string) error {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("id = ?", userID).Delete(&domain.User{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("delete user fail. error : %s", err.Error())
	}

	return tx.Commit().Error
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
