package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Login    string `json:"login"`    // email or mobile
	Password string `json:"password"` // user password
}

type PanInput struct {
	Pan string `json:"pan"`
}

type Users struct {
	UsersId   uint      `gorm:"primaryKey" json:"users_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func HashPassword(password string) string {
	// Generate a bcrypt hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hashedPassword)
}

func GetUserByEmail(db *gorm.DB, email string) (bool, error) {
	var count int64
	err := db.Model(&Users{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserByMobile(db *gorm.DB, username string) (bool, error) {
	var count int64
	err := db.Model(&Users{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserByLogin(db *gorm.DB, login, loginType string) (*Users, error) {
	var user Users
	var err error
	switch loginType {
	case "email":
		err = db.Where("email = ?", login).First(&user).Error
	case "mobile":
		err = db.Where("username = ?", login).First(&user).Error
	default:
		return nil, fmt.Errorf("invalid login type")
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsSamePassword(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}
