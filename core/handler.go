package core

import (
	"errors"
	"golang-todo-api-tdd-ddd/handler"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitializeHandler(e *echo.Echo, db *gorm.DB) error {

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return errors.New("JWT_SECRET environment variable not set")
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "online <3"})
	})

	handler.InitializeAuthenHandler(e, db)
	handler.InitializeTodoHandler(e, db, secretKey)
	handler.InitializeUserHandler(e, db, secretKey)

	return nil
}
