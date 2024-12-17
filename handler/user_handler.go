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

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitializeUserHandler(engine helper.Engine) {

	userRepo := repository.NewUserRepository(engine.DB)
	userService := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	userGroup := engine.Echo.Group("/user")

	userGroup.Use(echojwt.WithConfig(echojwt.Config{SigningKey: []byte(engine.SecretKey)}))

	userGroup.GET("/all", userHandler.GetAllUser)
	userGroup.GET("/:userID", userHandler.GetUserByID)
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

	allUserDTO := []valueobject.GetUserVO{}

	err := handler.userService.GetAllUser(&allUserDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to get all user. error: %s.", err))
	}

	return c.JSON(http.StatusOK, allUserDTO)
}

func (handler *UserHandler) GetUserByID(c echo.Context) error {

	response := core.ApiRespose{}

	userID := c.Param("userID")

	getUserVO, err := handler.userService.GetUserByID(userID)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Error = ""
	response.Data = getUserVO

	return c.JSON(http.StatusOK, response)
}

func (handler *UserHandler) UpdateUser(c echo.Context) error {

	user := domain.User{}
	userDTO := valueobject.UpdateUserVO{}

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
