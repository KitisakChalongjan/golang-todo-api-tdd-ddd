package main

import (
	"golang-todo-api-tdd-ddd/core"
	"log"
	"time"

	// echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	log.Println("initializing the application...")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      10 * time.Second,
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "time out(5s)",
	}))
	// e.Use(echojwt.WithConfig(echojwt.Config{
	// 	SigningKey: []byte("golang-todo-api-tdd-ddd"),
	// }))

	db, err := core.ConnectPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	core.InitializeHandler(e, db)

	if err = e.Start(":1323"); err != nil {

		log.Fatal(err)
	}
}
