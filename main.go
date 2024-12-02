package main

import (
	"golang-todo-api-tdd-ddd/core"
	"log"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {

	log.Println("initializing the application...")

	e := echo.New()

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("golang-todo-api-tdd-ddd"),
	}))

	db, err := core.ConnectPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	core.InitializeHandler(e, db)

	if err = e.Start(":1323"); err != nil {

		log.Fatal(err)
	}

}
