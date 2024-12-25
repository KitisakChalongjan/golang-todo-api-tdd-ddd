package service

import (
	"fmt"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/valueobject"

	"github.com/golang-jwt/jwt/v5"
)

type TodoService struct {
	todoRepo repository.ITodoRepository
}

func NewTodoService(todoRepo repository.ITodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

func (todoService *TodoService) GetTodoByID(todoID string, accessToken *jwt.Token) (valueobject.GetTodoVO, error) {

	tokenUserId, err := accessToken.Claims.GetSubject()
	if err != nil {
		return valueobject.GetTodoVO{}, fmt.Errorf("fail to get userID from access token: %w", err)
	}

	todoVO, err := todoService.todoRepo.GetTodoByID(todoID, tokenUserId)
	if err != nil {
		return valueobject.GetTodoVO{}, err
	}

	return todoVO, nil
}

func (todoService *TodoService) GetTodosByUserID(userID string, accessToken *jwt.Token) ([]valueobject.GetTodoVO, error) {

	tokenUserId, err := accessToken.Claims.GetSubject()
	if err != nil {
		return []valueobject.GetTodoVO{}, fmt.Errorf("fail to get todoID from access token: %w", err)
	}

	allTodoVO, err := todoService.todoRepo.GetTodosByUserID(userID, tokenUserId)
	if err != nil {
		return []valueobject.GetTodoVO{}, err
	}

	return allTodoVO, nil
}

func (todoService *TodoService) CreateTodo(todoDTO valueobject.CreateTodoVO) (string, error) {

	todo, err := todoService.todoRepo.CreateTodo(todoDTO)
	if err != nil {
		return "", err
	}

	return todo.ID, nil
}

func (todoService *TodoService) UpdateTodo(updateTodoVO valueobject.UpdateTodoVO, accessToken *jwt.Token) (string, error) {

	tokenUserId, err := accessToken.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("fail to get userID from access token: %w", err)
	}

	todo, err := todoService.todoRepo.UpdateTodo(updateTodoVO, tokenUserId)
	if err != nil {
		return "", err
	}

	return todo.ID, nil
}

func (todoService *TodoService) DeleteTodo(todoID string, accessToken *jwt.Token) (string, error) {

	tokenUserId, err := accessToken.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("fail to get userID from access token: %w", err)
	}

	todo, err := todoService.todoRepo.DeleteTodo(todoID, tokenUserId)
	if err != nil {
		return "", err
	}

	return todo.ID, nil
}
