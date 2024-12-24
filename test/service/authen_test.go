package service_test

import (
	"errors"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/helper"
	"golang-todo-api-tdd-ddd/service"
	"golang-todo-api-tdd-ddd/valueobject"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserById(userID string, tokenUserID string) (valueobject.GetUserVO, error) {
	args := m.Called(userID, tokenUserID)
	return args.Get(0).(valueobject.GetUserVO), args.Error(1)
}
func (m *MockUserRepository) GetUserByCredential(signInVO valueobject.SignInVO) (valueobject.GetUserVO, error) {
	args := m.Called(signInVO)
	return args.Get(0).(valueobject.GetUserVO), args.Error(1)
}

func (m *MockUserRepository) CreateUser(signUpVO valueobject.SignUpVO) (domain.User, error) {
	args := m.Called(signUpVO)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(updateUserVO valueobject.UpdateUserVO, tokenUserID string) (domain.User, error) {
	args := m.Called(updateUserVO, tokenUserID)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(userID string, tokenUserID string) (domain.User, error) {
	args := m.Called(userID, tokenUserID)
	return args.Get(0).(domain.User), args.Error(1)
}

func TestSuccessSignUp1(t *testing.T) {

	log.Info("testing TestSuccessSignUp1()....")

	mockRepo := new(MockUserRepository)

	authenService := service.NewAuthenService(mockRepo)

	signUpVO := valueobject.SignUpVO{
		Name:     "test",
		Email:    "test@example.com",
		Roles:    []string{"admin"},
		Username: "username",
		Password: "password",
	}

	mockRepo.On("CreateUser", mock.Anything).Return(domain.User{ID: "xxxx-xxxx-xxxx-xxxx"}, nil)

	userID, err := authenService.SignUp(signUpVO)

	assert.NoError(t, err)
	assert.Equal(t, "xxxx-xxxx-xxxx-xxxx", userID)
}

func TestErrorSignUp1(t *testing.T) {

	log.Info("testing TestErrorSignUp1()....")

	mockRepo := new(MockUserRepository)

	authenService := service.NewAuthenService(mockRepo)

	signUpVO := valueobject.SignUpVO{
		Name:     "test",
		Email:    "test@example.com",
		Roles:    []string{"admin"},
		Username: "username",
		Password: "password",
	}

	mockRepo.On("CreateUser", mock.Anything).Return(domain.User{ID: "xxxx-xxxx-xxxx-xxxx"}, errors.New("SignUp error!"))

	userID, err := authenService.SignUp(signUpVO)

	assert.Error(t, err)
	assert.Equal(t, "", userID)
	assert.EqualError(t, errors.New("SignUp error!"), err.Error())

	mockRepo.AssertExpectations(t)
}

func TestSuccessSignIn1(t *testing.T) {

	log.Info("testing TestSuccessSignIn1()....")

	mockRepo := new(MockUserRepository)

	authenService := service.NewAuthenService(mockRepo)

	signInVO := valueobject.SignInVO{
		Username: "username",
		Password: "password",
	}

	secretKey := "golang-todo-api-tdd-ddd"

	mockRepo.On("GetUserByCredential", mock.Anything).Return(valueobject.GetUserVO{ID: "xxxx-xxxx-xxxx-xxxx", Name: "name", Email: "test@example.com", Roles: []string{"guest"}}, nil)

	accessTokenString, err := authenService.SignIn(signInVO, secretKey)

	tokenClaims, _ := helper.TokenClaimsFromAccessTokenString(accessTokenString, secretKey)

	sub, _ := tokenClaims.Claims.GetSubject()

	assert.NoError(t, err)
	assert.Equal(t, "xxxx-xxxx-xxxx-xxxx", sub)

	mockRepo.AssertExpectations(t)
}

func TestErrorSignIn1(t *testing.T) {

	log.Info("testing TestErrorSignIn1()....")

	mockRepo := new(MockUserRepository)

	authenService := service.NewAuthenService(mockRepo)

	signInVO := valueobject.SignInVO{
		Username: "username",
		Password: "password",
	}

	secretKey := "golang-todo-api-tdd-ddd"

	mockRepo.On("GetUserByCredential", mock.Anything).Return(valueobject.GetUserVO{}, errors.New("SignIn error!"))

	accessTokenString, err := authenService.SignIn(signInVO, secretKey)

	assert.Error(t, err)
	assert.Equal(t, "", accessTokenString)
	assert.EqualError(t, errors.New("SignIn error!"), err.Error())

	mockRepo.AssertExpectations(t)
}
