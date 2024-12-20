package handler

import (
	"fmt"
	"golang-todo-api-tdd-ddd/core"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/helper"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/service"
	"golang-todo-api-tdd-ddd/valueobject"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitializeTodoHandler(engine helper.Engine) {
	todoRepo := repository.NewTodoRepository(engine.DB)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := NewTodoHandler(todoService)

	todoGroup := engine.Echo.Group("/todo")

	todoGroup.Use(echojwt.WithConfig(echojwt.Config{SigningKey: []byte(engine.SecretKey)}))

	// todoGroup.GET("/all", todoHandler.GetAllTodo)
	todoGroup.GET("/:todoID", todoHandler.GetTodoByID)
	todoGroup.GET("/user/:userID", todoHandler.GetTodosByUserID)
	todoGroup.POST("/", todoHandler.CreateTodo)
	todoGroup.PUT("/", todoHandler.UpdateTodo)
	todoGroup.DELETE("/:todoID", todoHandler.DeleteTodo)
}

type TodoHandler struct {
	todoService *service.TodoService
}

func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// func (handler *TodoHandler) GetAllTodo(c echo.Context) error {

// 	allTodoDTO := []valueobject.GetTodoVO{}

// 	if err := handler.todoService.GetAllTodo(&allTodoDTO); err != nil {
// 		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to get all todo. error: %s.", err))
// 	}

// 	return c.JSON(http.StatusOK, allTodoDTO)
// }

func (handler *TodoHandler) GetTodoByID(c echo.Context) error {

	response := core.ApiRespose{}

	todoID := c.Param("todoID")
	accessToken := c.Get("user").(*jwt.Token)

	todoVO, err := handler.todoService.GetTodoByID(todoID, accessToken)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Error = ""
	response.Data = todoVO

	return c.JSON(http.StatusOK, response)
}

func (handler *TodoHandler) GetTodosByUserID(c echo.Context) error {

	allTodoDTO := []valueobject.GetTodoVO{}

	userID := c.Param("userID")

	if err := handler.todoService.GetTodosByUserID(&allTodoDTO, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to get all todo by userID. error: %s.", err))
	}

	return c.JSON(http.StatusOK, allTodoDTO)
}

func (handler *TodoHandler) CreateTodo(c echo.Context) error {

	todo := domain.Todo{}
	todoDTO := valueobject.CreateTodoVO{}

	if err := c.Bind(&todoDTO); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request. error: %s.", err))
	}

	if err := handler.todoService.CreateTodo(&todo, todoDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to create todo. error: %s.", err))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("todo '%s' created", todo.ID))
}

func (handler *TodoHandler) UpdateTodo(c echo.Context) error {

	todo := domain.Todo{}
	todoDTO := valueobject.UpdateTodoVO{}

	if err := c.Bind(&todoDTO); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request. error: %s.", err))
	}

	if err := handler.todoService.UpdateTodo(&todo, todoDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to create todo. error: %s.", err))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("user '%s' updated", todo.ID))
}

func (handler *TodoHandler) DeleteTodo(c echo.Context) error {

	todoID := c.Param("todoID")

	if err := handler.todoService.DeleteTodo(todoID); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to deleted todo. error: %s.", err))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("user '%s' deleted", todoID))
}
