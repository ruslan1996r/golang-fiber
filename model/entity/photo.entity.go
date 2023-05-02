package entity

import (
	"time"

	"gorm.io/gorm"
)

// Photo Если не существует Категории с указанным ID, будет ошибка создания
type Photo struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Image      string         `json:"image"`
	CategoryID uint           `json:"category_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}
