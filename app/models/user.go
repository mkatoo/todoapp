package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	Name           string `json:"name"`
	Email          string `json:"email" gorm:"unique"`
	HashedPassword []byte `json:"-"`
	Tasks          []Task `json:"tasks" gorm:"foreignKey:UserID"`
	Token          Token  `json:"token" gorm:"foreignKey:UserID"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = hash
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
	return err == nil
}

func NewUser(name, email, password string) (*User, error) {
	user := &User{
		Name:  name,
		Email: email,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	return user, nil
}

func IsUserExists(db *gorm.DB, email string) (bool, error) {
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
