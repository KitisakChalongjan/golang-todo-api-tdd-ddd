package handler

import (
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitializeUserHandler(e *echo.Echo, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	userGroup := e.Group("/user")

	userGroup.GET("/all", func(c echo.Context) error {
		return GetAllUser(c, userService)
	})
	userGroup.GET("/:userID", func(c echo.Context) error {
		return GetUserByID(c, userService)
	})
	userGroup.POST("/", func(c echo.Context) error {
		return CreateUser(c, userService)
	})
	userGroup.PUT("/", func(c echo.Context) error {
		return UpdateUser(c, userService)
	})
	userGroup.DELETE("/:userID", func(c echo.Context) error {
		return DeleteUser(c, userService)
	})
}

func GetAllUser(c echo.Context, userService *service.UserService) error {

	users := []domain.User{}

	err := userService.GetAllUser(&users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func GetUserByID(c echo.Context, userService *service.UserService) error {

	user := domain.User{}

	userID := c.Param("userID")

	err := userService.GetUser(&user, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func CreateUser(c echo.Context, userService *service.UserService) error {

	user := domain.User{}
	userDTO := domain.CreateUserDTO{}

	err := c.Bind(&userDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request.")
	}

	err = userService.CreateUser(&user, &userDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "fail to create user.")
	}

	return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context, userService *service.UserService) error {

	user := domain.User{}
	userDTO := domain.UpdateUserDTO{}

	err := c.Bind(&userDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request :"+err.Error())
	}

	err = userService.UpdateUser(&user, &userDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "fail to update user : "+err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context, userService *service.UserService) error {

	userID := c.Param("userID")

	err := userService.DeleteUser(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "fail to delete user.")
	}

	return c.JSON(http.StatusOK, "delete user successful.")
}
