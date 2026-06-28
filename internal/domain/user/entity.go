package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	RoleDriver = "driver"
	RoleAdmin  = "admin"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(100);not null"`
	Role     string `gorm:"type:varchar(20);not null;default:driver"`
}

func (u *User) hashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) checkPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
