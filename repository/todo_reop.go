package repository

import (
	"golang-todo-api-tdd-ddd/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITodoRepository interface {
	GetAllTodos() (*[]domain.Todo, error)
	GetTodoByID(todoID uuid.UUID) (*domain.Todo, error)
	GetTodosByUserID(userID uuid.UUID) (*[]domain.Todo, error)
	CreateTodo(todo *domain.Todo) (*domain.Todo, error)
	UpdateTodo(todo *domain.Todo) (*domain.Todo, error)
	DeleteTodo(todoID uuid.UUID) error
}

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (todoRepo *TodoRepository) GetAllTodos() (*[]domain.Todo, error) {
	var todos []domain.Todo

	result := todoRepo.db.Find(&todos)

	return &todos, result.Error
}

func (todoRepo *TodoRepository) GetTodoByID(todoID uuid.UUID) (*domain.Todo, error) {
	var todo domain.Todo

	result := todoRepo.db.First(&todo, todoID)

	return &todo, result.Error
}

func (todoRepo *TodoRepository) GetTodosByUserID(userID uuid.UUID) (*[]domain.Todo, error) {
	var todos *[]domain.Todo

	result := todoRepo.db.Find(&todos, domain.Todo{UserID: userID})

	return todos, result.Error
}

func (todoRepo *TodoRepository) CreateTodo(todo *domain.Todo) (*domain.Todo, error) {
	result := todoRepo.db.Create(todo)

	return todo, result.Error
}

func (todoRepo *TodoRepository) UpdateTodo(todo *domain.Todo) (*domain.Todo, error) {
	result := todoRepo.db.Save(todo)

	return todo, result.Error
}

func (todoRepo *TodoRepository) DeleteTodo(todoID uuid.UUID) error {
	result := todoRepo.db.Delete(&domain.Todo{}, todoID)

	return result.Error
}
