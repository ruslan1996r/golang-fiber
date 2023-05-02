package handler

import (
	"fmt"
	"log"

	"fiber/database"
	"fiber/model/entity"
	"fiber/model/request"
	"fiber/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func PhotoHandlerCreate(ctx *fiber.Ctx) error {
	photo := new(request.PhotoCrateRequest)

	if err := ctx.BodyParser(photo); err != nil {
		return err
	}

	// Validate request
	validate := validator.New()
	errValidate := validate.Struct(photo)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	// Handle file, Validation required image
	filenames := ctx.Locals(utils.FILENAMES)

	if filenames == nil {
		return ctx.Status(422).JSON(fiber.Map{
			"message": "Images are required",
		})
	} else {
		// Saving files paths to database
		for _, filename := range filenames.([]string) { // filenames AS Array<string>
			newPhoto := entity.Photo{
				Image:      filename,
				CategoryID: photo.CategoryId,
			}

			errCreatePhoto := database.DB.Create(&newPhoto).Error
			if errCreatePhoto != nil {
				log.Println("Photo creation error", errCreatePhoto)
			}
		}
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message":   "success",
		"filenames": filenames,
	})
}

func PhotoHandlerDelete(ctx *fiber.Ctx) error {
	photoId := ctx.Params("id")

	var photo entity.Photo
	err := database.DB.First(&photo, "id = ?", photoId).Error

	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "photo not found",
		})
	}

	// Handle remove photo
	errDeleteFile := utils.HandleRemoveFile(photo.Image, "")

	if errDeleteFile != nil {
		fmt.Println("Fail to delete some file")
	}

	errDelete := database.DB.Delete(&photo).Error
	if errDelete != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "photo deletion from DB error",
		})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "Photo with id '" + photoId + "' was deleted",
	})
}
