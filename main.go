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

	_, err := os.Stat(".env")
	if err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading .env file")
		}
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

	if err = e.Start("0.0.0.0:1323"); err != nil {

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

// docker build -t todo-api
// docker run --name todo-api-1 -e JWT_SECRET=golang-todo-api-tdd-ddd -e DB_HOST=172.17.0.2 -e DB_PORT=5432 -e DB_USER=postgres -e DB_PASSWORD=Dewsmaller1* -e DB_NAME=postgres -p 1323:1323 -d todo-api
// docker run -d -v postgresql-data:/var/lib/postgresql/data -e POSTGRES_PASSWORD=Dewsmaller1\ -p 5433:5432 --name postgresql-main postgres
