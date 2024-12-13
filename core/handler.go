package core

import (
	"golang-todo-api-tdd-ddd/handler"
	"golang-todo-api-tdd-ddd/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitializeHandler(engine helper.Engine) error {

	engine.Echo.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "online <3"})
	})

	handler.InitializeAuthenHandler(engine)
	handler.InitializeTodoHandler(engine)
	handler.InitializeUserHandler(engine)

	return nil
}
