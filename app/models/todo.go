package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	*gorm.Model
	UserID  uint   `json:"user_id"`
	Content string `json:"content"`
}
