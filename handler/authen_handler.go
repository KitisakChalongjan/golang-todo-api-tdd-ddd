package handler

import (
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitializeAuthenHandler(e *echo.Echo, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	authenService := service.NewAuthenService(userRepo)
	authenHandler := NewAuthenHandler(authenService)

	authenGroup := e.Group("/authen")

	authenGroup.POST("/login", authenHandler.Login)
}

type AuthenHandler struct {
	authenService *service.AuthenService
}

func NewAuthenHandler(authenService *service.AuthenService) *AuthenHandler {
	return &AuthenHandler{authenService: authenService}
}

func (handler *AuthenHandler) Login(c echo.Context) error {

	var tokenString string
	loginDTO := domain.LogingDTO{}

	err := c.Bind(&loginDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request. error: %s.", err))
	}

	err = handler.authenService.Login(&tokenString, loginDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to login. error: %s.", err))
	}

	return c.JSON(http.StatusOK, map[string]string{"jwt": tokenString})
}
