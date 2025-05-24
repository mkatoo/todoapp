package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	*gorm.Model
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

func FindOrCreateToken(db *gorm.DB, userID uint) (*Token, error) {
	var existingToken Token
	if err := db.Where("user_id = ?", userID).First(&existingToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			newToken := &Token{
				UserID: userID,
				Token:  uuid.NewString(),
			}
			if err := db.Create(newToken).Error; err != nil {
				return nil, err
			}
			return newToken, nil
		}
		return nil, err
	}
	return &existingToken, nil
}
