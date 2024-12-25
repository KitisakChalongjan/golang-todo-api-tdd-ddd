package service_test

import (
	"errors"
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

func TestSuccessGetUserByID(t *testing.T) {

	t.Log("testing TestSuccessGetUserByID()....")

	mockRepo := new(MockUserRepository)

	userService := service.NewUserService(mockRepo)

	userID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(valueobject.GetUserVO{ID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", Name: "name", Username: "username", Email: "test@example.com", Roles: []string{"guest"}}, nil)

	getUserVO, err := userService.GetUserByID(userID, &accessToken)

	a1 := assert.NoError(t, err)
	a2 := assert.Equal(t, userID, getUserVO.ID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		t.Log("TestSuccessGetUserByID passed")
	}

	t.Log("TestSuccessGetUserByID failed")
}

func TestFailGetUserByID1(t *testing.T) {

	t.Log("testing TestFailGetUserByID1()....")

	mockRepo := new(MockUserRepository)

	userService := service.NewUserService(mockRepo)

	userID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(valueobject.GetUserVO{}, errors.New("GetUserByID error!"))

	getUserVO, err := userService.GetUserByID(userID, &accessToken)

	a1 := assert.Error(t, err)
	a2 := assert.Equal(t, valueobject.GetUserVO{}, getUserVO)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		t.Log("TestFailGetUserByID1 passed")
	}

	t.Log("TestFailGetUserByID1 failed")
}

func TestSuccessUpdateUser1(t *testing.T) {

	t.Log("testing TestSuccessUpdateUser1()....")

	mockRepo := new(MockUserRepository)

	userService := service.NewUserService(mockRepo)

	updateUserVO := valueobject.UpdateUserVO{
		ID:       "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		Name:     "newname",
		Email:    "test@example.com",
		Roles:    []string{"guest"},
		Password: "newpassword",
	}
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("UpdateUser", mock.Anything, mock.Anything).Return(domain.User{ID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", Name: "newname", Email: "test@example.com", Username: "username"}, nil)

	userID, err := userService.UpdateUser(updateUserVO, &accessToken)

	a1 := assert.NoError(t, err)
	a2 := assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", userID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		t.Log("TestSuccessUpdateUser1 passed")
	}

	t.Log("TestSuccessUpdateUser1 failed")
}

func TestFailUpdateUser1(t *testing.T) {

	t.Log("testing TestFailUpdateUser1()....")

	mockRepo := new(MockUserRepository)

	userService := service.NewUserService(mockRepo)

	updateUserVO := valueobject.UpdateUserVO{
		ID:       "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		Name:     "newname",
		Email:    "test@example.com",
		Roles:    []string{"guest"},
		Password: "newpassword",
	}
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("UpdateUser", mock.Anything, mock.Anything).Return(domain.User{}, errors.New("UpdateUser error!"))

	userID, err := userService.UpdateUser(updateUserVO, &accessToken)

	a1 := assert.Error(t, err)
	a2 := assert.Equal(t, "", userID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		t.Log("TestFailUpdateUser1 passed")
	}

	t.Log("TestFailUpdateUser1 failed")
}

func TestSuccessDeleteUser1(t *testing.T) {

	t.Log("testing TestSuccessDeleteUser1()....")

	mockRepo := new(MockUserRepository)

	userService := service.NewUserService(mockRepo)

	userID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("DeleteUser", mock.Anything, mock.Anything).Return(domain.User{ID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", Name: "name", Username: "username", Email: "test@example.com"}, nil)

	userID, err := userService.DeleteUser(userID, &accessToken)

	a1 := assert.NoError(t, err)
	a2 := assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", userID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		t.Log("TestSuccessDeleteUser1 passed")
	}

	t.Log("TestSuccessDeleteUser1 failed")
}

func TestFailDeleteUser1(t *testing.T) {

	t.Log("testing TestFailDeleteUser1()....")

	mockRepo := new(MockUserRepository)

	userService := service.NewUserService(mockRepo)

	userID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("DeleteUser", mock.Anything, mock.Anything).Return(domain.User{}, errors.New("DeleteUser error!"))

	userID, err := userService.DeleteUser(userID, &accessToken)

	a1 := assert.Error(t, err)
	a2 := assert.Equal(t, "", userID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		t.Log("TestFailDeleteUser1 passed")
	}

	t.Log("TestFailDeleteUser1 failed")
}
