package handler

import (
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitializeTodoHandler(e *echo.Echo, db *gorm.DB) {
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)

	todoGroup := e.Group("/todo")

	todoGroup.GET("/all", func(c echo.Context) error {
		return GetAllTodo(c, todoService)
	})
	todoGroup.GET("/:todoID", func(c echo.Context) error {
		return GetTodoByID(c, todoService)
	})
	todoGroup.GET("/user/:userID", func(c echo.Context) error {
		return GetTodosByUserID(c, todoService)
	})
	todoGroup.POST("/", func(c echo.Context) error {
		return CreateTodo(c, todoService)
	})
	todoGroup.PUT("/", func(c echo.Context) error {
		return UpdateTodo(c, todoService)
	})
	todoGroup.DELETE("/:todoID", func(c echo.Context) error {
		return DeleteTodo(c, todoService)
	})
}

func GetAllTodo(c echo.Context, todoService *service.TodoService) error {

	todos := []domain.Todo{}

	err := todoService.GetAllTodo(&todos)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, todos)
}

func GetTodoByID(c echo.Context, todoService *service.TodoService) error {

	todo := domain.Todo{}

	todoID := c.Param("todoID")

	err := todoService.GetTodoByID(&todo, todoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, todo)
}

func GetTodosByUserID(c echo.Context, todoService *service.TodoService) error {

	todos := []domain.Todo{}

	userID := c.Param("userID")

	err := todoService.GetTodosByUserID(&todos, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, todos)
}

func CreateTodo(c echo.Context, todoService *service.TodoService) error {

	todo := domain.Todo{}
	todoDTO := domain.CreateTodoDTO{}

	err := c.Bind(&todoDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request.")
	}

	err = todoService.CreateTodo(&todo, todoDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "fail to create user.")
	}

	return c.JSON(http.StatusOK, todo)
}

func UpdateTodo(c echo.Context, todoService *service.TodoService) error {

	todo := domain.Todo{}
	todoDTO := domain.UpdateTodoDTO{}

	err := c.Bind(&todoDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request :"+err.Error())
	}

	err = todoService.UpdateTodo(&todo, todoDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "fail to update todo : "+err.Error())
	}

	return c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c echo.Context, todoService *service.TodoService) error {

	todoID := c.Param("todoID")

	err := todoService.DeleteTodo(todoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "fail to delete todo.")
	}

	return c.JSON(http.StatusOK, "delete todo successful.")
}
