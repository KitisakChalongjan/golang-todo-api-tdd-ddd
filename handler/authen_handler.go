package handler

import (
	"fmt"
	"golang-todo-api-tdd-ddd/core"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/helper"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/service"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitializeAuthenHandler(engine helper.Engine) {
	userRepo := repository.NewUserRepository(engine.DB)
	authenRepo := repository.NewAuthenRepository(engine.DB)
	authenService := service.NewAuthenService(userRepo, authenRepo)
	authenHandler := NewAuthenHandler(authenService)

	authenGroup := engine.Echo.Group("/authen")

	authenGroup.POST("/signup", authenHandler.SignUp)
	authenGroup.POST("/signin", authenHandler.SignIn)
	authenGroup.POST("/logout", authenHandler.Logout, echojwt.WithConfig(echojwt.Config{SigningKey: []byte(engine.SecretKey)}))
	authenGroup.POST("/re-access-token", authenHandler.ReAccessToken)
}

type AuthenHandler struct {
	authenService *service.AuthenService
}

func NewAuthenHandler(authenService *service.AuthenService) *AuthenHandler {
	return &AuthenHandler{authenService: authenService}
}

func (handler *AuthenHandler) SignUp(c echo.Context) error {

	response := core.ApiRespose{}
	signupDTO := domain.SignUpDTO{}

	err := c.Bind(&signupDTO)
	if err != nil {
		response.Error = err.Error()
		response.Data = nil
		return c.JSON(http.StatusBadRequest, response)
	}

	userID, err := handler.authenService.SignUpUser(signupDTO)
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
	loginDTO := domain.SignInDTO{}

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

func (handler *AuthenHandler) Logout(c echo.Context) error {

	accessToken := c.Get("user").(*jwt.Token)

	if err := handler.authenService.Logout(accessToken); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to logout. error: %s.", err))
	}

	return c.JSON(http.StatusOK, "logout success")
}

func (handler *AuthenHandler) ReAccessToken(c echo.Context) error {

	var newAccessToken string
	var newRefreshToken string
	reAccessDTO := domain.ReAccessDTO{}

	reAccessDTO.DeviceInfo = c.Request().Header.Get("User-Agent")
	reAccessDTO.IpAddress = c.RealIP()

	refreshTokenCookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	err = handler.authenService.ReAccessToken(refreshTokenCookie, &newAccessToken, &newRefreshToken, reAccessDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	return c.JSON(http.StatusOK, map[string]string{"access_token": newAccessToken})
}
