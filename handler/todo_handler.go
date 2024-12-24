package handler

import (
	"fmt"
	"golang-todo-api-tdd-ddd/core"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/service"
	"golang-todo-api-tdd-ddd/valueobject"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitializeTodoHandler(engine core.Engine) {
	todoRepo := repository.NewTodoRepository(engine.DB)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := NewTodoHandler(todoService)

	todoGroup := engine.Echo.Group("/todo")

	todoGroup.Use(echojwt.WithConfig(echojwt.Config{SigningKey: []byte(engine.SecretKey)}))

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

	response := core.ApiRespose{}

	userID := c.Param("userID")
	accessToken := c.Get("user").(*jwt.Token)

	allTodoVO, err := handler.todoService.GetTodosByUserID(userID, accessToken)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Error = ""
	response.Data = allTodoVO

	return c.JSON(http.StatusOK, allTodoVO)
}

func (handler *TodoHandler) CreateTodo(c echo.Context) error {

	response := core.ApiRespose{}
	todoDTO := valueobject.CreateTodoVO{}

	err := c.Bind(&todoDTO)
	if err != nil {
		response.Error = fmt.Errorf("invalid request: %s", err).Error()
		response.Data = nil
		return c.JSON(http.StatusBadRequest, response)
	}

	todoID, err := handler.todoService.CreateTodo(todoDTO)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Error = ""
	response.Data = map[string]string{"message": fmt.Errorf("todo '%s' created", todoID).Error()}

	return c.JSON(http.StatusOK, response)
}

func (handler *TodoHandler) UpdateTodo(c echo.Context) error {

	response := core.ApiRespose{}
	updateTodoVO := valueobject.UpdateTodoVO{}

	accessToken := c.Get("user").(*jwt.Token)

	if err := c.Bind(&updateTodoVO); err != nil {
		response.Error = fmt.Errorf("invalid request: %s", err).Error()
		response.Data = nil
		return c.JSON(http.StatusBadRequest, response)
	}

	todoID, err := handler.todoService.UpdateTodo(updateTodoVO, accessToken)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Error = ""
	response.Data = map[string]string{"message": fmt.Errorf("todo '%s' updated", todoID).Error()}

	return c.JSON(http.StatusOK, response)
}

func (handler *TodoHandler) DeleteTodo(c echo.Context) error {

	response := core.ApiRespose{}

	todoID := c.Param("todoID")
	accessToken := c.Get("user").(*jwt.Token)

	todoID, err := handler.todoService.DeleteTodo(todoID, accessToken)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Error = ""
	response.Data = map[string]string{"message": fmt.Errorf("todo '%s' deleted", todoID).Error()}

	return c.JSON(http.StatusOK, response)
}
