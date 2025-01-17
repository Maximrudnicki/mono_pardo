package users

import (
	"errors"
	"regexp"
	"strings"

	"mono_pardo/internal/utils"
)

type User struct {
	Id       int    `gorm:"type:int;primary_key"`
	Username string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

func NewUser(username, email, password string) (*User, error) {
	hashedPassword, err := utils.HashPassword(strings.TrimSpace(password))
	if err != nil {
		return nil, err
	}

	validUsername := strings.TrimSpace(username)
	validEmail := strings.TrimSpace(email)

	if !isValidEmail(validEmail) {
		return nil, errors.New("invalid email format")
	}

	return &User{
		Username: validUsername,
		Email:    validEmail,
		Password: hashedPassword,
	}, nil
}

func isValidEmail(email string) bool {
	// Regular expression for validating an Email
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
