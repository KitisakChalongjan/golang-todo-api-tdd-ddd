package repository

import (
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/valueobject"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITodoRepository interface {
	GetTodoByID(todoID string, tokenUserID string) (valueobject.GetTodoVO, error)
	GetTodosByUserID(userID string, tokenUserID string) ([]valueobject.GetTodoVO, error)
	CreateTodo(createTodoVO valueobject.CreateTodoVO) (domain.Todo, error)
	UpdateTodo(updateTodoVO valueobject.UpdateTodoVO, tokenUserID string) (domain.Todo, error)
	DeleteTodo(todoID string) error
}

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (repo *TodoRepository) GetTodoByID(todoID string, tokenUserID string) (valueobject.GetTodoVO, error) {

	todo := domain.Todo{}
	todoVO := valueobject.GetTodoVO{}

	err := repo.db.Where("id = ?", todoID).First(&todo).Error
	if err != nil {
		return valueobject.GetTodoVO{}, fmt.Errorf("faile to get todo by id: %w", err)
	}

	if tokenUserID != todo.UserID {
		return valueobject.GetTodoVO{}, fmt.Errorf("you cannot access this data: userID not match")
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

func (repo *TodoRepository) GetTodosByUserID(userID string, tokenUserID string) ([]valueobject.GetTodoVO, error) {

	allTodo := []domain.Todo{}

	err := repo.db.Where("user_id = ?", userID).Find(&allTodo).Error
	if err != nil {
		return []valueobject.GetTodoVO{}, fmt.Errorf("faile to get all todo by userID: %w", err)
	}

	if tokenUserID != allTodo[0].UserID {
		return []valueobject.GetTodoVO{}, fmt.Errorf("you cannot access this data: userID not match")
	}

	allTodoVO := make([]valueobject.GetTodoVO, len(allTodo))

	for i, todo := range allTodo {
		allTodoVO[i] = valueobject.GetTodoVO{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Title,
			IsCompleted: todo.IsCompleted,
			Priority:    todo.Priority,
			Due:         todo.Due,
			UserID:      todo.UserID,
		}
	}

	return allTodoVO, nil
}

func (repo *TodoRepository) CreateTodo(createTodoVO valueobject.CreateTodoVO) (domain.Todo, error) {

	todo := domain.Todo{}

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Error
	if err != nil {
		return domain.Todo{}, err
	}

	todo.ID = uuid.New().String()
	todo.Title = createTodoVO.Title
	todo.Description = createTodoVO.Description
	todo.IsCompleted = false
	todo.Priority = createTodoVO.Priority
	todo.Due = createTodoVO.Due
	todo.UserID = createTodoVO.UserID

	if err := repo.db.Create(&todo).Error; err != nil {
		tx.Rollback()
		return domain.Todo{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.Todo{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return todo, nil
}

func (repo *TodoRepository) UpdateTodo(updateTodoVO valueobject.UpdateTodoVO, tokenUserID string) (domain.Todo, error) {

	todo := domain.Todo{}

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Error
	if err != nil {
		return domain.Todo{}, fmt.Errorf("transaction fail: %s", err.Error())
	}

	err = tx.Where("id = ?", updateTodoVO.ID).First(&todo).Error
	if err != nil {
		tx.Rollback()
		return domain.Todo{}, fmt.Errorf("no todo found: %s", err.Error())
	}

	if tokenUserID != todo.UserID {
		return domain.Todo{}, fmt.Errorf("you cannot access this data: userID not match")
	}

	todo.Title = updateTodoVO.Title
	todo.Description = updateTodoVO.Description
	todo.IsCompleted = updateTodoVO.IsCompleted
	todo.Priority = updateTodoVO.Priority
	todo.Due = updateTodoVO.Due

	if err := tx.Save(&todo).Error; err != nil {
		tx.Rollback()
		return domain.Todo{}, fmt.Errorf("update todo fail: %s", err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.Todo{}, fmt.Errorf("update todo commit fail: %s", err.Error())
	}

	return todo, nil
}

func (repo *TodoRepository) DeleteTodo(todoID string, tokenUserID string) (domain.Todo, error) {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return domain.Todo{}, err
	}

	todo := domain.Todo{}

	err := tx.Where("id = ?", todoID).First(&todo).Error
	if err != nil {
		tx.Rollback()
		return domain.Todo{}, fmt.Errorf("no todo found: %s", err.Error())
	}

	if tokenUserID != todo.UserID {
		return domain.Todo{}, fmt.Errorf("you cannot access this data: userID not match")
	}

	if err := repo.db.Where("id = ?", todoID).Delete(&domain.Todo{}).Error; err != nil {
		tx.Rollback()
		return domain.Todo{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.Todo{}, fmt.Errorf("delete todo commit fail: %s", err.Error())
	}

	return todo, nil
}
