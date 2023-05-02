package handler

import (
	"fmt"

	"fiber/database"
	"fiber/model/entity"
	"fiber/model/request"
	"fiber/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func BookHandlerCreate(ctx *fiber.Ctx) error {
	book := new(request.BookCreateRequest)

	if err := ctx.BodyParser(book); err != nil {
		return err
	}

	// Validate request
	validate := validator.New()
	errValidate := validate.Struct(book)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	// Handle file, Validation required image
	filename := ctx.Locals(utils.FILENAME)
	if filename == nil {
		return ctx.Status(422).JSON(fiber.Map{
			"message": "Image cover is required",
		})
	}

	filenameString := fmt.Sprintf("%v", filename) // To string

	newBook := entity.Book{
		Title:  book.Title,
		Author: book.Author,
		Cover:  filenameString,
	}

	errCreateUser := database.DB.Create(&newBook).Error
	if errCreateUser != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to save book",
			"error":   errCreateUser,
		})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "success",
		"data":    newBook,
	})
}
