package main

import (
	"fmt"

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

	err := app.Listen(":8080")

	if err != nil {
		fmt.Println("Server Init Error!")
	}
}
