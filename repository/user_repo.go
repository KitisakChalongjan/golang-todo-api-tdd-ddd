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
	GetAllUsers(users *[]domain.User) error
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

func (repo *UserRepository) GetAllUsers(allUserDTO *[]valueobject.GetUserVO) error {

	allUser := []domain.User{}

	if err := repo.db.Find(&allUser).Error; err != nil {
		return fmt.Errorf("no users found. error : %s", err.Error())
	}

	*allUserDTO = make([]valueobject.GetUserVO, len(allUser))

	for i, user := range allUser {
		(*allUserDTO)[i] = valueobject.GetUserVO{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email,
			ProfileImgURL: user.ProfileImgURL,
			Username:      user.Username,
			CreatedAt:     user.CreatedAt,
		}
	}

	return nil
}

func (repo *UserRepository) GetUserById(userID string) (valueobject.GetUserVO, error) {

	user := domain.User{}
	userVO := valueobject.GetUserVO{}

	if err := repo.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return valueobject.GetUserVO{}, err
	}

	roleNames, err := GetRoleNamesByUserID(repo, user.ID)
	if err != nil {
		return valueobject.GetUserVO{}, err
	}

	// usersRoles := []domain.UsersRoles{}

	// err := repo.db.Where("user_id = ?", user.ID).Find(&usersRoles).Error
	// if err != nil {
	// 	return valueobject.GetUserVO{}, err
	// }

	// roleIDs := []string{}

	// for _, usersRolesElement := range usersRoles {
	// 	roleIDs = append(roleIDs, usersRolesElement.RoleID)
	// }

	// roles := []domain.Role{}

	// err = repo.db.Where("id IN (?)", roleIDs).Find(&roles).Error
	// if err != nil {
	// 	return valueobject.GetUserVO{}, err
	// }

	// roleNames := []string{}

	// for _, role := range roles {
	// 	roleNames = append(roleNames, role.Name)
	// }

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
		return valueobject.GetUserVO{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(signInVO.Password)); err != nil {
		return valueobject.GetUserVO{}, err
	}

	roleNames, err := GetRoleNamesByUserID(repo, user.ID)
	if err != nil {
		return valueobject.GetUserVO{}, err
	}

	// usersRoles := []domain.UsersRoles{}

	// err := repo.db.Where("user_id = ?", user.ID).Find(&usersRoles).Error
	// if err != nil {
	// 	return valueobject.GetUserVO{}, err
	// }

	// roleIDs := []string{}

	// for _, usersRolesElement := range usersRoles {
	// 	roleIDs = append(roleIDs, usersRolesElement.RoleID)
	// }

	// roles := []domain.Role{}

	// err = repo.db.Where("id IN (?)", roleIDs).Find(&roles).Error
	// if err != nil {
	// 	return valueobject.GetUserVO{}, err
	// }

	// roleNames := []string{}

	// for _, role := range roles {
	// 	roleNames = append(roleNames, role.Name)
	// }

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

	// create new user
	newUser := domain.User{}

	newUser.ID = uuid.New().String()
	newUser.Name = signupDTO.Name
	newUser.Email = signupDTO.Email
	newUser.ProfileImgURL = signupDTO.ProfileImgURL
	newUser.Username = signupDTO.Username
	newUser.PasswordHash = signupDTO.Password

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return domain.User{}, err
	}

	// get roles id by roles name
	roles := []domain.Role{}

	err := repo.db.Where("name IN (?)", signupDTO.Roles).Find(&roles).Error
	if err != nil {
		return domain.User{}, err
	}

	roleIDs := []string{}

	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	// create new users roles
	newUserRoles := []domain.UsersRoles{}

	for _, roleID := range roleIDs {
		newUserRole := domain.UsersRoles{
			UserID: newUser.ID,
			RoleID: roleID,
		}
		newUserRoles = append(newUserRoles, newUserRole)
	}

	if err := tx.Create(&newUserRoles).Error; err != nil {
		tx.Rollback()
		return domain.User{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}

func (repo *UserRepository) UpdateUser(user *domain.User, userDTO *valueobject.UpdateUserVO) error {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("id = ?", userDTO.ID).First(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("no user found. error : %s", err.Error())
	}

	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.ProfileImgURL = userDTO.ProfileImgURL
	user.PasswordHash = userDTO.Password

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("update user fail. error : %s", err.Error())
	}

	return tx.Commit().Error
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
		return []string{}, err
	}

	roleIDs := []string{}

	for _, usersRolesElement := range usersRoles {
		roleIDs = append(roleIDs, usersRolesElement.RoleID)
	}

	roles := []domain.Role{}

	err = repo.db.Where("id IN (?)", roleIDs).Find(&roles).Error
	if err != nil {
		return []string{}, err
	}

	roleNames := []string{}

	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	return roleNames, nil
}
