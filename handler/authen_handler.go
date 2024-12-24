package handler

import (
	"fmt"
	"golang-todo-api-tdd-ddd/core"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/service"
	"golang-todo-api-tdd-ddd/valueobject"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func InitializeAuthenHandler(engine core.Engine) {
	userRepo := repository.NewUserRepository(engine.DB)
	// roleRepo := repository.NewRoleRepository(engine.DB)
	// userRoleRepo := repository.NewUserRoleRepository(engine.DB)
	authenService := service.NewAuthenService(userRepo /*, roleRepo, userRoleRepo*/)
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

	response.Error = ""
	response.Data = map[string]string{"message": fmt.Sprintf("user '%s' created", userID)}

	return c.JSON(http.StatusOK, response)
}

func (handler *AuthenHandler) SignIn(c echo.Context) error {

	response := core.ApiRespose{}
	signInVO := valueobject.SignInVO{}

	err := c.Bind(&signInVO)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusBadRequest, response)
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return fmt.Errorf("JWT_SECRET environment variable not set")
	}

	accessTokenString, err := handler.authenService.SignIn(signInVO, secretKey)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Data = map[string]string{"access_token": accessTokenString}

	return c.JSON(http.StatusOK, response)
}
