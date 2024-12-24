package service_test

import (
	"errors"
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/service"
	"golang-todo-api-tdd-ddd/valueobject"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(userID string, tokenUserID string) (valueobject.GetUserVO, error) {
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

func TestFailGetUserByID1(t *testing.T) {

	fmt.Println("testing TestFailGetUserByID1()....")

	mockRepo := new(MockUserRepository)

	userService := service.NewUserService(mockRepo)

	userID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(valueobject.GetUserVO{}, errors.New("GetUserByID error!"))

	getUserVO, err := userService.GetUserByID(userID, &accessToken)

	assert.Error(t, err)
	assert.Equal(t, valueobject.GetUserVO{}, getUserVO)

	mockRepo.AssertExpectations(t)
}
