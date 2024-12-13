package handler

import (
	"fmt"
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

	authenGroup.POST("/signup", authenHandler.SignUpUser)
	authenGroup.POST("/login", authenHandler.Login)
	authenGroup.POST("/logout", authenHandler.Logout, echojwt.WithConfig(echojwt.Config{SigningKey: []byte(engine.SecretKey)}))
	authenGroup.POST("/re-access-token", authenHandler.ReAccessToken, echojwt.WithConfig(echojwt.Config{SigningKey: []byte(engine.SecretKey)}))
}

type AuthenHandler struct {
	authenService *service.AuthenService
}

func NewAuthenHandler(authenService *service.AuthenService) *AuthenHandler {
	return &AuthenHandler{authenService: authenService}
}

func (handler *AuthenHandler) SignUpUser(c echo.Context) error {

	user := domain.User{}
	userDTO := domain.SignUpUserDTO{}

	if err := c.Bind(&userDTO); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request. error: %s.", err))
	}

	if err := handler.authenService.SignUpUser(&user, userDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to create user. error: %s.", err))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("user '%s' created", user.ID))
}

func (handler *AuthenHandler) Login(c echo.Context) error {

	var accessToken string
	var refreshToken string
	loginDTO := domain.LoginDTO{}

	if err := c.Bind(&loginDTO); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid request. error: %s.", err))
	}

	loginDTO.DeviceInfo = c.Request().Header.Get("User-Agent")
	loginDTO.IpAddress = c.RealIP()

	err := handler.authenService.Login(&accessToken, &refreshToken, loginDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to login. error: %s.", err))
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Use true if you are using HTTPS
	})

	return c.JSON(http.StatusOK, map[string]string{"access_token": accessToken})
}

func (handler *AuthenHandler) Logout(c echo.Context) error {

	accessToken := c.Get("user").(*jwt.Token)
	fmt.Println(accessToken)

	if err := handler.authenService.Logout(accessToken); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("fail to logout. error: %s.", err))
	}

	return c.JSON(http.StatusOK, "")
}

func (handler *AuthenHandler) ReAccessToken(c echo.Context) error {

	// if refresh_token, err := c.Cookie("refresh_token"); err != nil {
	// 	return c.JSON(http.StatusUnauthorized, "refresh token is not found.")
	// }
	return nil
}
