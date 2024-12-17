package core

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ApiRespose struct {
	Error string      `json:"message"`
	Data  interface{} `json:"data"`
}

type Engine struct {
	Echo      *echo.Echo
	DB        *gorm.DB
	SecretKey string
}
