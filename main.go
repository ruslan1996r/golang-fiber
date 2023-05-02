package main

import (
	"fmt"
	"os"

	"fiber/database"
	"fiber/database/migration"
	"fiber/route"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.DBInit()
	migration.RunMigration()

	config := fiber.Config{
		Concurrency: 100000,
	}

	app := fiber.New(config)

	route.RouterInit(app)

	appErr := app.Listen(":8080")

	if appErr != nil {
		fmt.Println("Server Init Error!")
		os.Exit(1)
	}
}
