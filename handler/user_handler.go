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

func InitializeUserHandler(engine core.Engine) {

	userRepo := repository.NewUserRepository(engine.DB)
	userService := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	userGroup := engine.Echo.Group("/user")

	userGroup.Use(echojwt.WithConfig(echojwt.Config{SigningKey: []byte(engine.SecretKey)}))

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

func (handler *UserHandler) GetUserByID(c echo.Context) error {

	response := core.ApiRespose{}

	userID := c.Param("userID")
	accessToken := c.Get("user").(*jwt.Token)

	getUserVO, err := handler.userService.GetUserByID(userID, accessToken)
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

	response := core.ApiRespose{}
	updateUserDTO := valueobject.UpdateUserVO{}

	accessToken := c.Get("user").(*jwt.Token)

	if err := c.Bind(&updateUserDTO); err != nil {
		response.Error = fmt.Errorf("invalid request: %s", err).Error()
		response.Data = nil
		return c.JSON(http.StatusBadRequest, response)
	}

	userID, err := handler.userService.UpdateUser(updateUserDTO, accessToken)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Error = ""
	response.Data = map[string]string{"message": fmt.Errorf("user '%s' updated", userID).Error()}

	return c.JSON(http.StatusOK, response)
}

func (handler *UserHandler) DeleteUser(c echo.Context) error {

	response := core.ApiRespose{}

	userID := c.Param("userID")
	accessToken := c.Get("user").(*jwt.Token)

	userID, err := handler.userService.DeleteUser(userID, accessToken)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Error = ""
	response.Data = map[string]string{"message": fmt.Sprintf("user '%s' deleted", userID)}

	return c.JSON(http.StatusOK, response)
}
