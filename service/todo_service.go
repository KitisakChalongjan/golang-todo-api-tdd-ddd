package service

import (
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
)

type TodoService struct {
	todoRepo *repository.TodoRepository
}

func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

func (todoService *TodoService) GetAllTodo(allTodoDTO *[]domain.GetTodoDTO) error {

	err := todoService.todoRepo.GetAllTodos(allTodoDTO)
	if err != nil {
		return err
	}

	return nil
}

func (todoService *TodoService) GetTodoByID(todo *domain.GetTodoDTO, todoID string) error {

	if err := todoService.todoRepo.GetTodoByID(todo, todoID); err != nil {
		return err
	}

	return nil
}

func (todoService *TodoService) GetTodosByUserID(allTodoDTO *[]domain.GetTodoDTO, userID string) error {

	result := todoService.todoRepo.GetTodosByUserID(allTodoDTO, userID)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (todoService *TodoService) CreateTodo(todo *domain.Todo, todoDTO domain.CreateTodoDTO) error {

	if err := todoService.todoRepo.CreateTodo(todo, todoDTO); err != nil {
		return err
	}

	return nil
}

func (todoService *TodoService) UpdateTodo(todo *domain.Todo, todoDTO domain.UpdateTodoDTO) error {

	if err := todoService.todoRepo.UpdateTodo(todo, todoDTO); err != nil {
		return err
	}

	return nil
}

func (todoService *TodoService) DeleteTodo(todoID string) error {

	if err := todoService.todoRepo.DeleteTodo(todoID); err != nil {
		return err
	}

	return nil
}
