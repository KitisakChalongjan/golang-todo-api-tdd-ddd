package main

import (
	"golang-todo-api-tdd-ddd/core"
	"golang-todo-api-tdd-ddd/handler"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	log.Println("initializing the application...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	db, err := core.ConnectPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      5 * time.Second,
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "time out(5s)",
	}))

	engine := core.Engine{
		Echo:      e,
		DB:        db,
		SecretKey: secretKey,
	}

	if err := InitializeHandler(engine); err != nil {
		log.Fatal(err)
	}

	if err = e.Start(":1323"); err != nil {

		log.Fatal(err)
	}
}

func InitializeHandler(engine core.Engine) error {

	engine.Echo.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "online <3"})
	})

	handler.InitializeRoleHandler(engine)
	handler.InitializeAuthenHandler(engine)
	handler.InitializeTodoHandler(engine)
	handler.InitializeUserHandler(engine)

	return nil
}
