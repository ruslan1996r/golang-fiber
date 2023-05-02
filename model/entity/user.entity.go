package entity

import (
	"time"

	"gorm.io/gorm"
)

// User json "-" означает, что это поле скрытое для юзера и его нельзя получить при GET-запросе
// Но поля можно получить при SQL-запросе
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"index:unique"`
	Email     string         `json:"email"`
	Password  string         `json:"-" gorm:"column:password"`
	Address   string         `json:"address"`
	Phone     string         `json:"phone"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}
