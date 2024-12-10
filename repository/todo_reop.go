package repository

import (
	"golang-todo-api-tdd-ddd/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITodoRepository interface {
	GetAllTodos(todos *[]domain.Todo) error
	GetTodoByID(todo *domain.Todo, todoID string) error
	GetTodosByUserID(todos *[]domain.Todo, userID string) *gorm.DB
	CreateTodo(todo *domain.Todo, todoDTO domain.CreateTodoDTO) error
	UpdateTodo(todo *domain.Todo, todoDTO domain.UpdateTodoDTO) error
	DeleteTodo(todoID string) error
}

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (repo *TodoRepository) GetAllTodos(allTodoDTO *[]domain.GetTodoDTO) error {

	allTodo := []domain.Todo{}

	if err := repo.db.Find(&allTodo).Error; err != nil {
		return err
	}

	*allTodoDTO = make([]domain.GetTodoDTO, len(allTodo))

	for i, todo := range allTodo {
		(*allTodoDTO)[i] = domain.GetTodoDTO{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Title,
			IsCompleted: todo.IsCompleted,
			Priority:    todo.Priority,
			Due:         todo.Due,
			UserID:      todo.UserID,
		}
	}

	return nil
}

func (repo *TodoRepository) GetTodoByID(todoDTO *domain.GetTodoDTO, todoID string) error {

	todo := domain.Todo{}

	if err := repo.db.Where("id = ?", todoID).First(&todo).Error; err != nil {
		return err
	}

	todoDTO.ID = todo.ID
	todoDTO.Title = todo.Title
	todoDTO.Description = todo.Description
	todoDTO.IsCompleted = todo.IsCompleted
	todoDTO.Priority = todo.Priority
	todoDTO.Due = todo.Due
	todoDTO.UserID = todo.UserID

	return nil
}

func (repo *TodoRepository) GetTodosByUserID(allTodoDTO *[]domain.GetTodoDTO, userID string) *gorm.DB {

	allTodo := []domain.Todo{}

	result := repo.db.Where("user_id = ?", userID).Find(&allTodo)

	*allTodoDTO = make([]domain.GetTodoDTO, len(allTodo))

	for i, todo := range allTodo {
		(*allTodoDTO)[i] = domain.GetTodoDTO{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Title,
			IsCompleted: todo.IsCompleted,
			Priority:    todo.Priority,
			Due:         todo.Due,
			UserID:      todo.UserID,
		}
	}

	return result
}

func (repo *TodoRepository) CreateTodo(todo *domain.Todo, todoDTO domain.CreateTodoDTO) error {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	todo.ID = uuid.New().String()
	todo.Title = todoDTO.Title
	todo.Description = todoDTO.Description
	todo.IsCompleted = false
	todo.Priority = todoDTO.Priority
	todo.Due = todoDTO.Due
	todo.UserID = todoDTO.UserID

	if err := repo.db.Create(&todo).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *TodoRepository) UpdateTodo(todo *domain.Todo, todoDTO domain.UpdateTodoDTO) error {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("id = ?", todoDTO.ID).First(todo).Error; err != nil {
		tx.Rollback()
		return err
	}

	todo.Title = todoDTO.Title
	todo.Description = todoDTO.Description
	todo.IsCompleted = todoDTO.IsCompleted
	todo.Priority = todoDTO.Priority
	todo.Due = todoDTO.Due

	if err := tx.Save(&todo).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *TodoRepository) DeleteTodo(todoID string) error {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := repo.db.Where("id = ?", todoID).Delete(&domain.Todo{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
