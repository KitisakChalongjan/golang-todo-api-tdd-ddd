package repository

import (
	"golang-todo-api-tdd-ddd/domain"

	"gorm.io/gorm"
)

type ITodoRepository interface {
	GetAllTodos(todos *[]domain.Todo) *gorm.DB
	GetTodoByID(todo *domain.Todo, todoID string) *gorm.DB
	GetTodosByUserID(todos *[]domain.Todo, userID string) *gorm.DB
	CreateTodo(todo *domain.Todo) *gorm.DB
	UpdateTodo(todo *domain.Todo) *gorm.DB
	DeleteTodo(todoID string) error
}

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (todoRepo *TodoRepository) GetAllTodos(todos *[]domain.Todo) *gorm.DB {

	result := todoRepo.db.Find(&todos)

	return result
}

func (todoRepo *TodoRepository) GetTodoByID(todo *domain.Todo, todoID string) *gorm.DB {

	result := todoRepo.db.Where("id = ?", todoID).First(&todo)

	return result
}

func (todoRepo *TodoRepository) GetTodosByUserID(todos *[]domain.Todo, userID string) *gorm.DB {

	result := todoRepo.db.Where("user_id = ?", userID).Find(&todos)

	return result
}

func (todoRepo *TodoRepository) CreateTodo(todo *domain.Todo) *gorm.DB {

	result := todoRepo.db.Create(&todo)

	return result
}

func (todoRepo *TodoRepository) UpdateTodo(todo *domain.Todo) *gorm.DB {

	result := todoRepo.db.Save(&todo)

	return result
}

func (todoRepo *TodoRepository) DeleteTodo(todoID string) *gorm.DB {

	result := todoRepo.db.Where("id = ?", todoID).Delete(&domain.Todo{})

	return result
}
