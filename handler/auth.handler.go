package handler

import (
	"fmt"
	"time"

	"fiber/database"
	"fiber/model/entity"
	"fiber/model/request"
	"fiber/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func LoginHandler(ctx *fiber.Ctx) error {
	loginRequest := new(request.LoginRequest)

	if err := ctx.BodyParser(loginRequest); err != nil {
		return err
	}

	// Validate request
	validate := validator.New()
	errValidate := validate.Struct(loginRequest)
	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	var user entity.User
	err := database.DB.First(&user, "email = ?", loginRequest.Email).Error
	if err != nil {
		return ctx.Status(402).JSON(fiber.Map{
			"message": "User does not exists",
		})
	}

	// Check Validation Password
	isValid := utils.CheckPasswordHash(loginRequest.Password, user.Password)
	if !isValid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Wrong credentials",
		})
	}

	// Generate token
	claims := jwt.MapClaims{
		// "name": user.Name,
		// "email":user.Email,
		// "address": user.Address,
		// "exp": time.Now().Add(time.Minute * 2).Unix(),
	}
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["address"] = user.Address
	claims["exp"] = time.Now().Add(time.Minute * 2).Unix()

	if user.Email == "zhora@gmail.com" {
		claims["role"] = "admin"
	} else {
		claims["role"] = "user"
	}

	token, errGenerateToken := utils.GenerateToken(&claims)

	if errGenerateToken != nil {
		fmt.Println("errGenerateToken", errGenerateToken)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token error",
		})
	}

	return ctx.Status(400).JSON(fiber.Map{
		"message": "success",
		"token":   token,
	})
}
