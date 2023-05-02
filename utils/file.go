package utils

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
)

const PATH_TO_STATIC = "./public/covers/"

var AllowedTypes = []string{"image/jpg", "image/jpeg"}

func HandleSingleFile(ctx *fiber.Ctx) error {
	file, errFile := ctx.FormFile(COVER)
	if errFile != nil {
		log.Println("Error File = ", errFile)
	}

	var filename *string
	if file != nil {
		errCheckContentType := checkContentType(file, "image/jpg", "image/jpeg") // if !Contains(AllowedTypes, contentTypeFile) {
		if errCheckContentType != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": errCheckContentType.Error(),
			})
		}

		filename = &file.Filename
		// Создание кастомного имени файла
		// extensionFile := filepath.Ext(*filename)
		// newFileName := fmt.Sprintf("%s%s", uuid(), extensionFile)

		errSaveFile := ctx.SaveFile(file, fmt.Sprintf("./public/covers/%s", *filename)) // newFileName
		if errSaveFile != nil {
			log.Println("Fail to store file")
		}
	} else {
		log.Println("Nothing file to uploading")
	}

	ctx.Locals(FILENAME, *filename)

	return ctx.Next()
}

func HandleMultipleFiles(ctx *fiber.Ctx) error {
	form, errForm := ctx.MultipartForm()
	if errForm != nil {
		log.Println("Error Read Multipart from Request, Error = ", errForm)
	}

	files := form.File["photos"]

	var filenames []string
	for i, file := range files {
		var filename string
		if file != nil {
			filename = fmt.Sprintf("%d_%s", i, file.Filename)

			errSaveFile := ctx.SaveFile(file, fmt.Sprintf("./public/covers/%s", filename))
			if errSaveFile != nil {
				log.Println("Fail to store file")
			}
		} else {
			log.Println("Nothing file to uploading")
		}

		if filename != "" {
			filenames = append(filenames, filename)
		}
	}

	ctx.Locals(FILENAMES, filenames)

	return ctx.Next()
}

func HandleRemoveFile(filename, pathFile string) error {
	if len(pathFile) > 0 {
		err := os.Remove(pathFile + filename)
		if err != nil {
			log.Println("Failed to remove file")
			return err
		}
	} else {
		err := os.Remove(PATH_TO_STATIC + filename)
		if err != nil {
			log.Println("Failed to remove file")
			return err
		}
	}

	return nil
}

func checkContentType(file *multipart.FileHeader, contentTypes ...string) error {
	if len(contentTypes) > 0 {
		for _, contentType := range contentTypes {
			contentTypeFile := file.Header.Get("Content-Type")
			if contentTypeFile == contentType {
				return nil
			}
		}

		return errors.New("not allowed file type")
	}

	return errors.New("not found content type to be checking")
}
