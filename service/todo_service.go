package service

import (
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"

	"github.com/google/uuid"
)

type TodoService struct {
	todoRepo *repository.TodoRepository
}

func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

func (todoService *TodoService) GetAllTodo(todos *[]domain.Todo) error {

	result := todoService.todoRepo.GetAllTodos(todos)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (todoService *TodoService) GetTodoByID(todo *domain.Todo, todoID string) error {

	result := todoService.todoRepo.GetTodoByID(todo, todoID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (todoService *TodoService) GetTodosByUserID(todos *[]domain.Todo, userID string) error {

	result := todoService.todoRepo.GetTodosByUserID(todos, userID)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (todoService *TodoService) CreateTodo(todo *domain.Todo, todoDTO domain.CreateTodoDTO) error {

	todo.ID = uuid.New().String()
	todo.Title = todoDTO.Title
	todo.Description = todoDTO.Description
	todo.IsCompleted = false
	todo.Priority = todoDTO.Priority
	todo.Due = todoDTO.Due
	todo.UserID = todoDTO.UserID

	result := todoService.todoRepo.CreateTodo(todo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (todoService *TodoService) UpdateTodo(todo *domain.Todo, todoDTO domain.UpdateTodoDTO) error {

	result := todoService.todoRepo.GetTodoByID(todo, todoDTO.ID)
	if result.Error != nil {
		return result.Error
	}

	todo.Title = todoDTO.Title
	todo.Description = todoDTO.Description
	todo.Priority = todoDTO.Priority
	todo.Due = todoDTO.Due

	result = todoService.todoRepo.UpdateTodo(todo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (todoService *TodoService) DeleteTodo(todoID string) error {

	result := todoService.todoRepo.DeleteTodo(todoID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
