package handler

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"fiber/database"
	"fiber/model/entity"
	"fiber/model/request"
	"fiber/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UserHandlerGetAll(ctx *fiber.Ctx) error {
	// userInfo := ctx.Locals("userInfo") // Получаю эту инфу из Мидлвара
	// log.Println("userInfo", userInfo)

	var users []entity.User

	result := database.DB.Find(&users)

	if result.Error != nil {
		log.Println(result.Error)
	}

	return ctx.Status(200).JSON(users)
}

// UserHandlerCreate new() - что-то вроде make(), но в отличии от него возвращает ссылку, а не сам объект
func UserHandlerCreate(ctx *fiber.Ctx) error {
	user := new(request.UserCreateReq)

	if err := ctx.BodyParser(user); err != nil {
		return err
	}

	// Validate request
	validate := validator.New()
	errValidate := validate.Struct(user)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	newUser := entity.User{
		Name:     user.Name,
		Email:    user.Email,
		Address:  user.Address,
		Phone:    user.Phone,
		Password: user.Password,
	}

	hashedPassword, err := utils.HashingPassword(user.Password)

	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	newUser.Password = hashedPassword

	errCreateUser := database.DB.Create(&newUser).Error
	if errCreateUser != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to store data",
			"error":   errCreateUser,
		})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "success",
		"data":    newUser,
	})
}

func UserHandlerGetById(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")

	var user entity.User
	// var user response.UserResponse
	// err := database.DB.Model(entity.User{ ID: userId }).First(&user).Error  // Можно делать запрос так
	err := database.DB.First(&user, "id = ?", userId).Error

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "user with ID: '" + userId + "' not found",
		})
	}

	return ctx.Status(400).JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})
}

func UserHandlerUpdate(ctx *fiber.Ctx) error {
	userRequest := new(request.UserUpdateReq)
	if err := ctx.BodyParser(userRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "bad request",
		})
	}

	var user entity.User

	// Check Available User
	userId := ctx.Params("id")
	err := database.DB.First(&user, "id = ?", userId).Error
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "user with ID: '" + userId + "' not found",
		})
	}

	// Update User Data
	if userRequest.Name != "" {
		user.Name = userRequest.Name
	}
	user.Address = userRequest.Address
	user.Phone = userRequest.Phone

	errUpdate := database.DB.Save(&user).Error

	if errUpdate != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "User update error",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})
}

func UserHandlerUpdateEmail(ctx *fiber.Ctx) error {
	userRequest := new(request.UserEmailReq)
	if err := ctx.BodyParser(userRequest); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "bad request",
		})
	}

	// Validate request
	validate := validator.New()
	errValidate := validate.Struct(userRequest)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	var user entity.User
	var isEmailUserExists entity.User

	// Check Available User
	userId := ctx.Params("id")
	err := database.DB.First(&user, "id = ?", userId).Error
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "user with ID: '" + userId + "' not found",
		})
	}

	errCheckEmail := database.DB.First(&isEmailUserExists, "email = ?", userRequest.Email).Error
	if errCheckEmail == nil {
		return ctx.Status(402).JSON(fiber.Map{
			"message": "Email already in use",
		})
	}

	// Update User Data
	user.Email = userRequest.Email

	errUpdate := database.DB.Save(&user).Error

	if errUpdate != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "User update error",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})
}

func UserHandlerDelete(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")

	var user entity.User
	err := database.DB.First(&user, "id = ?", userId).Error

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "user with ID: '" + userId + "' not found",
		})
	}

	errDelete := database.DB.Delete(&user).Error
	if errDelete != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Error for user deleting",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "user with ID: '" + userId + "' was deleted",
	})
}

// DoSmt --------------------------------------------------------------------------------
func DoSmt(wg *sync.WaitGroup, counter *uint64) {
	fmt.Println("START...")

	for i := 0; i < 1000000; i++ {
		for i := 0; i < 1000; i++ {
			atomic.AddUint64(counter, 1)
		}
	}

	defer wg.Done()
	defer fmt.Println("END!!!!")
}

func TestMultithreading(ctx *fiber.Ctx) error {
	var wg sync.WaitGroup
	var counter uint64

	wg.Add(1)
	go DoSmt(&wg, &counter)

	wg.Wait()

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Success!",
		"index":   counter,
	})
}
