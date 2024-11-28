package main

import (
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	echo := echo.New()

	db, err := ConnectPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(db.Error)

	err = echo.Start(":1323")

	if err != nil {
		log.Fatal(err)
	}
}
