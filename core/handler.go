package core

import (
	"golang-todo-api-tdd-ddd/handler"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitializeHandler(e *echo.Echo, db *gorm.DB) {

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "online <3"})
	})

	handler.InitializeTodoHandler(e, db)
	handler.InitializeUserHandler(e, db)
}
