package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashingPassword(password string) (string, error) {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}

	return string(hashedByte), err
}

// CheckPasswordHash Вернёт Bool, True - когда пароль валидный, False - невалидный
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	fmt.Println("CompareHashAndPassword", err)
	return err == nil
}
