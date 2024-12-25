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

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) GetTodoByID(todoID string, tokenUserID string) (valueobject.GetTodoVO, error) {
	args := m.Called(todoID, tokenUserID)
	return args.Get(0).(valueobject.GetTodoVO), args.Error(1)
}

func (m *MockTodoRepository) GetTodosByUserID(userID string, tokenUserID string) ([]valueobject.GetTodoVO, error) {
	args := m.Called(userID, tokenUserID)
	return args.Get(0).([]valueobject.GetTodoVO), args.Error(1)
}

func (m *MockTodoRepository) CreateTodo(createTodoVO valueobject.CreateTodoVO) (domain.Todo, error) {
	args := m.Called(createTodoVO)
	return args.Get(0).(domain.Todo), args.Error(1)
}

func (m *MockTodoRepository) UpdateTodo(updateTodoVO valueobject.UpdateTodoVO, tokenUserID string) (domain.Todo, error) {
	args := m.Called(updateTodoVO, tokenUserID)
	return args.Get(0).(domain.Todo), args.Error(1)
}

func (m *MockTodoRepository) DeleteTodo(todoID string, tokenUserID string) (domain.Todo, error) {
	args := m.Called(todoID)
	return args.Get(0).(domain.Todo), args.Error(1)
}

func TestSuccessGetTodoByID(t *testing.T) {

	t.Log("testing TestSuccessGetTodoByID()....")

	mockRepo := new(MockTodoRepository)

	todoService := service.NewTodoService(mockRepo)

	todoID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("GetTodoByID", mock.Anything, mock.Anything).Return(valueobject.GetTodoVO{ID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", Title: "title", Description: "description", IsCompleted: false, Priority: "low", UserID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}, nil)

	getTodoVO, err := todoService.GetTodoByID(todoID, &accessToken)

	a1 := assert.NoError(t, err)
	a2 := assert.Equal(t, todoID, getTodoVO.ID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		t.Log("TestSuccessGetTodoByID passed")
	}

	t.Log("TestSuccessGetTodoByID failed")
}

func TestFailGetTodoByID(t *testing.T) {

	t.Log("testing TestFailGetTodoByID....")

	mockRepo := new(MockTodoRepository)

	todoService := service.NewTodoService(mockRepo)

	todoID := "11111111-1111-1111-1111-111111111111"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("GetTodoByID", mock.Anything, mock.Anything).Return(valueobject.GetTodoVO{}, errors.New("GetTodoByID error"))

	getTodoVO, err := todoService.GetTodoByID(todoID, &accessToken)

	a1 := assert.Error(t, err)
	a2 := assert.Equal(t, "", getTodoVO.ID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		t.Log("TestFailGetTodoByID passed")
	}

	t.Log("TestFailGetTodoByID failed")
}

func TestSuccessGetTodosByUserID(t *testing.T) {

	t.Log("testing TestSuccessGetTodosByUserID....")

	mockRepo := new(MockTodoRepository)

	todoService := service.NewTodoService(mockRepo)

	userID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("GetTodosByUserID", mock.Anything, mock.Anything).Return([]valueobject.GetTodoVO{{ID: "11111111-1111-1111-1111-111111111111", UserID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}, {ID: "22222222-2222-2222-2222-222222222222", UserID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}, nil)

	getTodoVO, err := todoService.GetTodosByUserID(userID, &accessToken)

	a1 := assert.NoError(t, err)
	a2 := assert.Equal(t, userID, getTodoVO[0].UserID)
	a3 := assert.Equal(t, userID, getTodoVO[len(getTodoVO)-1].UserID)
	a4 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 && a4 {
		t.Log("TestSuccessGetTodosByUserID passed")
	}

	t.Log("TestSuccessGetTodosByUserID failed")
}

func TestFailGetTodosByUserID(t *testing.T) {

	t.Log("testing TestFailGetTodosByUserID....")

	mockRepo := new(MockTodoRepository)

	todoService := service.NewTodoService(mockRepo)

	userID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("GetTodosByUserID", mock.Anything, mock.Anything).Return([]valueobject.GetTodoVO{}, errors.New("GetTodosByUserID error"))

	getTodoVO, err := todoService.GetTodosByUserID(userID, &accessToken)

	a1 := assert.Error(t, err)
	a2 := assert.Equal(t, 0, len(getTodoVO))
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		t.Log("TestFailGetTodosByUserID passed")
	}

	t.Log("TestFailGetTodosByUserID failed")
}
