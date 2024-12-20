package repository

import (
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/valueobject"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITodoRepository interface {
	// GetAllTodos(todos *[]domain.Todo) error
	GetTodoByID(todo *domain.Todo, todoID string) error
	GetTodosByUserID(todos *[]domain.Todo, userID string) *gorm.DB
	CreateTodo(todo *domain.Todo, todoDTO valueobject.CreateTodoVO) error
	UpdateTodo(todo *domain.Todo, todoDTO valueobject.UpdateTodoVO) error
	DeleteTodo(todoID string) error
}

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// func (repo *TodoRepository) GetAllTodos(allTodoDTO *[]valueobject.GetTodoVO) error {

// 	allTodo := []domain.Todo{}

// 	if err := repo.db.Find(&allTodo).Error; err != nil {
// 		return err
// 	}

// 	*allTodoDTO = make([]valueobject.GetTodoVO, len(allTodo))

// 	for i, todo := range allTodo {
// 		(*allTodoDTO)[i] = valueobject.GetTodoVO{
// 			ID:          todo.ID,
// 			Title:       todo.Title,
// 			Description: todo.Title,
// 			IsCompleted: todo.IsCompleted,
// 			Priority:    todo.Priority,
// 			Due:         todo.Due,
// 			UserID:      todo.UserID,
// 		}
// 	}

// 	return nil
// }

func (repo *TodoRepository) GetTodoByID(todoID string) (valueobject.GetTodoVO, error) {

	todo := domain.Todo{}
	todoVO := valueobject.GetTodoVO{}

	err := repo.db.Where("id = ?", todoID).First(&todo).Error
	if err != nil {
		return valueobject.GetTodoVO{}, fmt.Errorf("faile to get todo by id: %w", err)
	}

	todoVO.ID = todo.ID
	todoVO.Title = todo.Title
	todoVO.Description = todo.Description
	todoVO.IsCompleted = todo.IsCompleted
	todoVO.Priority = todo.Priority
	todoVO.Due = todo.Due
	todoVO.UserID = todo.UserID

	return todoVO, nil
}

func (repo *TodoRepository) GetTodosByUserID(allTodoDTO *[]valueobject.GetTodoVO, userID string) *gorm.DB {

	allTodo := []domain.Todo{}

	result := repo.db.Where("user_id = ?", userID).Find(&allTodo)

	*allTodoDTO = make([]valueobject.GetTodoVO, len(allTodo))

	for i, todo := range allTodo {
		(*allTodoDTO)[i] = valueobject.GetTodoVO{
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

func (repo *TodoRepository) CreateTodo(todo *domain.Todo, todoDTO valueobject.CreateTodoVO) error {

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

func (repo *TodoRepository) UpdateTodo(todo *domain.Todo, todoDTO valueobject.UpdateTodoVO) error {

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
