package service

import (
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"

	"github.com/google/uuid"
)

type TodoService struct {
	todoRepo repository.TodoRepository
}

func NewTodoService(todoRepo repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

func (todoService *TodoService) GetTodosByUserID(userId uuid.UUID) (*[]domain.Todo, error) {

	todos, err := todoService.todoRepo.GetTodosByUserID(userId)
	if err != nil {

	}

	return todos, nil
}
