package service_test

import (
	"errors"
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/service"
	"golang-todo-api-tdd-ddd/valueobject"
	"testing"
	"time"

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

	fmt.Println("testing TestSuccessGetTodoByID()....")

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
		fmt.Println("TestSuccessGetTodoByID passed")
		fmt.Println(" ")
		return
	}

	fmt.Println("TestSuccessGetTodoByID failed")
}

func TestFailGetTodoByID(t *testing.T) {

	fmt.Println("testing TestFailGetTodoByID....")

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
		fmt.Println("TestFailGetTodoByID passed")
		fmt.Println(" ")
		return
	}

	fmt.Println("TestFailGetTodoByID failed")
}

func TestSuccessGetTodosByUserID(t *testing.T) {

	fmt.Println("testing TestSuccessGetTodosByUserID....")

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
		fmt.Println("TestSuccessGetTodosByUserID passed")
		fmt.Println(" ")
		return
	}

	fmt.Println("TestSuccessGetTodosByUserID failed")
}

func TestFailGetTodosByUserID(t *testing.T) {

	fmt.Println("testing TestFailGetTodosByUserID....")

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
		fmt.Println("TestFailGetTodosByUserID passed")
		fmt.Println(" ")
		return
	}

	fmt.Println("TestFailGetTodosByUserID failed")
}

func TestSuccessCreateTodo(t *testing.T) {

	fmt.Println("testing TestSuccessCreateTodo....")

	mockRepo := new(MockTodoRepository)

	todoService := service.NewTodoService(mockRepo)

	due := time.Now().Add(time.Hour * 3)

	createTodoVO := valueobject.CreateTodoVO{
		Title:       "title",
		Description: "description",
		Priority:    "high",
		Due:         &due,
	}

	mockRepo.On("CreateTodo", mock.Anything).Return(domain.Todo{ID: "11111111-1111-1111-1111-111111111111", UserID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}, nil)

	todoID, err := todoService.CreateTodo(createTodoVO)

	a1 := assert.NoError(t, err)
	a2 := assert.Equal(t, "11111111-1111-1111-1111-111111111111", todoID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		fmt.Println("TestSuccessCreateTodo passed")
		fmt.Println(" ")
		return
	}

	fmt.Println("TestSuccessCreateTodo failed")
}

func TestFailCreateTodo(t *testing.T) {

	fmt.Println("testing TestFailCreateTodo....")

	mockRepo := new(MockTodoRepository)

	todoService := service.NewTodoService(mockRepo)

	due := time.Now().Add(time.Hour * 6)

	createTodoVO := valueobject.CreateTodoVO{
		Title:       "title",
		Description: "description",
		Priority:    "high",
		Due:         &due,
	}

	mockRepo.On("CreateTodo", mock.Anything).Return(domain.Todo{}, errors.New("CreateTodo error"))

	todoID, err := todoService.CreateTodo(createTodoVO)

	a1 := assert.Error(t, err)
	a2 := assert.Equal(t, "", todoID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		fmt.Println("TestFailCreateTodo passed")
		fmt.Println(" ")
		return
	}

	fmt.Println("TestFailCreateTodo failed")
}

func TestSuccessUpdateTodo(t *testing.T) {

	fmt.Println("testing TestSuccessUpdateTodo....")

	mockRepo := new(MockTodoRepository)

	todoService := service.NewTodoService(mockRepo)

	due := time.Now().Add(time.Hour * 3)

	updateTodoVO := valueobject.UpdateTodoVO{
		ID:          "11111111-1111-1111-1111-111111111111",
		Title:       "title",
		Description: "description",
		Priority:    "high",
		IsCompleted: true,
		Due:         &due,
	}
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("UpdateTodo", mock.Anything, mock.Anything).Return(domain.Todo{ID: "11111111-1111-1111-1111-111111111111", UserID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}, nil)

	todoID, err := todoService.UpdateTodo(updateTodoVO, &accessToken)

	a1 := assert.NoError(t, err)
	a2 := assert.Equal(t, updateTodoVO.ID, todoID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		fmt.Println("TestSuccessUpdateTodo passed")
		fmt.Println(" ")
		return
	}

	fmt.Println("TestSuccessUpdateTodo failed")
}

func TestFailUpdateTodo(t *testing.T) {

	fmt.Println("testing TestFailUpdateTodo....")

	mockRepo := new(MockTodoRepository)

	todoService := service.NewTodoService(mockRepo)

	due := time.Now().Add(time.Hour * 3)

	updateTodoVO := valueobject.UpdateTodoVO{
		ID:          "11111111-1111-1111-1111-111111111111",
		Title:       "title",
		Description: "description",
		Priority:    "high",
		IsCompleted: true,
		Due:         &due,
	}
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("UpdateTodo", mock.Anything, mock.Anything).Return(domain.Todo{}, errors.New("UpdateTodo error"))

	todoID, err := todoService.UpdateTodo(updateTodoVO, &accessToken)

	a1 := assert.Error(t, err)
	a2 := assert.Equal(t, "", todoID)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		fmt.Println("TestFailUpdateTodo passed")
		fmt.Println(" ")
		return
	}

	fmt.Println("TestFailUpdateTodo failed")
}

func TestSuccessDeleteTodo(t *testing.T) {

	fmt.Println("testing TestSuccessDeleteTodo....")

	mockRepo := new(MockTodoRepository)

	todoService := service.NewTodoService(mockRepo)

	todoID := "11111111-1111-1111-1111-111111111111"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("DeleteTodo", mock.Anything, mock.Anything).Return(domain.Todo{ID: "11111111-1111-1111-1111-111111111111", UserID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}, nil)

	todoIDResult, err := todoService.DeleteTodo(todoID, &accessToken)

	a1 := assert.NoError(t, err)
	a2 := assert.Equal(t, todoID, todoIDResult)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		fmt.Println("TestSuccessDeleteTodo passed")
		fmt.Println(" ")
		return
	}

	fmt.Println("TestSuccessDeleteTodo failed")
}

func TestFailDeleteTodo(t *testing.T) {

	fmt.Println("testing TestFailDeleteTodo....")

	mockRepo := new(MockTodoRepository)

	todoService := service.NewTodoService(mockRepo)

	todoID := "11111111-1111-1111-1111-111111111111"
	accessToken := jwt.Token{Claims: jwt.RegisteredClaims{Subject: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}

	mockRepo.On("DeleteTodo", mock.Anything, mock.Anything).Return(domain.Todo{}, errors.New("DeleteTodo error"))

	todoIDResult, err := todoService.DeleteTodo(todoID, &accessToken)

	a1 := assert.Error(t, err)
	a2 := assert.Equal(t, "", todoIDResult)
	a3 := mockRepo.AssertExpectations(t)

	if a1 && a2 && a3 {
		fmt.Println("TestFailDeleteTodo passed")
		fmt.Println(" ")
		return
	}

	fmt.Println("TestFailDeleteTodo failed")
}
