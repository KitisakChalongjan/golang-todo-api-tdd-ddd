package helper

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Engine struct {
	Echo      *echo.Echo
	DB        *gorm.DB
	SecretKey string
}
