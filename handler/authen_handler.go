package handler

import (
	"fmt"
	"golang-todo-api-tdd-ddd/core"
	"golang-todo-api-tdd-ddd/helper"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/service"
	"golang-todo-api-tdd-ddd/valueobject"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitializeAuthenHandler(engine helper.Engine) {
	userRepo := repository.NewUserRepository(engine.DB)
	roleRepo := repository.NewRoleRepository(engine.DB)
	userRoleRepo := repository.NewUserRoleRepository(engine.DB)
	authenService := service.NewAuthenService(userRepo, roleRepo, userRoleRepo)
	authenHandler := NewAuthenHandler(authenService)

	authenGroup := engine.Echo.Group("/authen")

	authenGroup.POST("/signup", authenHandler.SignUp)
	authenGroup.POST("/signin", authenHandler.SignIn)
}

type AuthenHandler struct {
	authenService *service.AuthenService
}

func NewAuthenHandler(authenService *service.AuthenService) *AuthenHandler {
	return &AuthenHandler{authenService: authenService}
}

func (handler *AuthenHandler) SignUp(c echo.Context) error {

	response := core.ApiRespose{}
	signupDTO := valueobject.SignUpVO{}

	err := c.Bind(&signupDTO)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusBadRequest, response)
	}

	userID, err := handler.authenService.SignUp(signupDTO)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Data = map[string]string{"message": fmt.Sprintf("create user id %s", userID)}

	return c.JSON(http.StatusOK, response)
}

func (handler *AuthenHandler) SignIn(c echo.Context) error {

	response := core.ApiRespose{}
	loginDTO := valueobject.SignInVO{}

	err := c.Bind(&loginDTO)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusBadRequest, response)
	}

	accessTokenString, err := handler.authenService.SignIn(loginDTO)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Data = map[string]string{"access_token": accessTokenString}

	return c.JSON(http.StatusOK, response)
}
