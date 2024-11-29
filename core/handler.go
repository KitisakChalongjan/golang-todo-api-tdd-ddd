package core

import (
	"golang-todo-api-tdd-ddd/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitializeHandler(e *echo.Echo, db *gorm.DB) {

	handler.InitializeTodoHandler(e, db)
	handler.InitializeUserHandler(e, db)
}
