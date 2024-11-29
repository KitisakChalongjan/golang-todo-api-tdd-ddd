package main

import (
	"golang-todo-api-tdd-ddd/core"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	db, err := core.ConnectPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Okay")
	})

	core.InitializeHandler(e, db)

	err = e.Start(":1323")

	if err != nil {
		log.Fatal(err)
	}
}
