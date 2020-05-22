package models

import (
	"server2/app/validator"

	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required"`

	PasswordHash string `json:"=" validate:"required"`
	Password     string `json:"password,omitempty" validate:"gt=5,lt=15"`
}

// ValidateBeforeStoring used before storing in database
func (u *User) ValidateBeforeStoring() error {
	return validator.Validate.Struct(*u)
}

// GeneratePasswordHash ...
func (u *User) GeneratePasswordHash() error {
	if u.Password == "" {
		return ErrPasswordIsEmpty
	}

	hash, err := generatePasswordHash(u.Password)

	if err != nil {
		return err
	}

	u.PasswordHash = hash
	return nil
}

// CheckPassword ...
func (u *User) CheckPassword(password string) error {
	if u.PasswordHash == "" {
		err := u.GeneratePasswordHash()
		if err != nil {
			return err
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	if err != nil {
		return err
	}

	return nil
}

func generatePasswordHash(pass string) (string, error) {
	if pass == "" {
		return "", ErrPasswordIsEmpty
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}
