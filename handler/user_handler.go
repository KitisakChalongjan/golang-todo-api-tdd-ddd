package handler

import (
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/service"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitializeUserHandler(e *echo.Echo, db *gorm.DB, secretKey string) {

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	userGroup := e.Group("/user")

	userGroup.Use(echojwt.WithConfig(echojwt.Config{SigningKey: []byte(secretKey)}))

	userGroup.GET("/all", userHandler.GetAllUser)
	userGroup.GET("/:userID", userHandler.GetUserByID)
	userGroup.POST("/", userHandler.CreateUser)
	userGroup.PUT("/", userHandler.UpdateUser)
	userGroup.DELETE("/:userID", userHandler.DeleteUser)
}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (handler *UserHandler) GetAllUser(c echo.Context) error {

	allUserDTO := []domain.GetUserDTO{}

	err := handler.userService.GetAllUser(&allUserDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to get all user. error: %s.", err))
	}

	return c.JSON(http.StatusOK, allUserDTO)
}

func (handler *UserHandler) GetUserByID(c echo.Context) error {

	user := domain.GetUserDTO{}

	userID := c.Param("userID")

	err := handler.userService.GetUser(&user, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to get user. error: %s.", err))
	}

	return c.JSON(http.StatusOK, user)
}

func (handler *UserHandler) CreateUser(c echo.Context) error {

	user := domain.User{}
	userDTO := domain.CreateUserDTO{}

	if err := c.Bind(&userDTO); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request. error: %s.", err))
	}

	if err := handler.userService.CreateUser(&user, userDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to create user. error: %s.", err))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("user '%s' created", user.ID))
}

func (handler *UserHandler) UpdateUser(c echo.Context) error {

	user := domain.User{}
	userDTO := domain.UpdateUserDTO{}

	if err := c.Bind(&userDTO); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request. error: %s.", err))
	}

	if err := handler.userService.UpdateUser(&user, &userDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to update user. error: %s.", err))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("user '%s' updated", user.ID))
}

func (handler *UserHandler) DeleteUser(c echo.Context) error {

	userID := c.Param("userID")

	if err := handler.userService.DeleteUser(userID); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("fail to delete user. error: %s.", err))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("user '%s' deleted", userID))
}
