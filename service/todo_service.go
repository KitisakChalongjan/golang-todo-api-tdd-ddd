package service

import (
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/valueobject"
)

type TodoService struct {
	todoRepo *repository.TodoRepository
}

func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

func (todoService *TodoService) GetAllTodo(allTodoDTO *[]valueobject.GetTodoVO) error {

	err := todoService.todoRepo.GetAllTodos(allTodoDTO)
	if err != nil {
		return err
	}

	return nil
}

func (todoService *TodoService) GetTodoByID(todo *valueobject.GetTodoVO, todoID string) error {

	if err := todoService.todoRepo.GetTodoByID(todo, todoID); err != nil {
		return err
	}

	return nil
}

func (todoService *TodoService) GetTodosByUserID(allTodoDTO *[]valueobject.GetTodoVO, userID string) error {

	result := todoService.todoRepo.GetTodosByUserID(allTodoDTO, userID)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (todoService *TodoService) CreateTodo(todo *domain.Todo, todoDTO valueobject.CreateTodoVO) error {

	if err := todoService.todoRepo.CreateTodo(todo, todoDTO); err != nil {
		return err
	}

	return nil
}

func (todoService *TodoService) UpdateTodo(todo *domain.Todo, todoDTO valueobject.UpdateTodoVO) error {

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
