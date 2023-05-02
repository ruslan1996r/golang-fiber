package migration

import (
	"fmt"
	"log"

	"fiber/database"
	"fiber/model/entity"
)

func RunMigration() {
	err := database.DB.AutoMigrate(
		&entity.User{},
		&entity.Book{},
		&entity.Category{},
		&entity.Photo{},
	)

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Database Migrated")
}
