package models

import (
	"github.com/niakr1s/chatty-server/app/er"
	"golang.org/x/crypto/bcrypt"
)

// Pass ...
type Pass struct {
	Password     string `json:"password,omitempty" validate:"gt=5,lt=15"`
	PasswordHash string `json:"-" validate:"required"`
}

// GeneratePasswordHash ...
func (p *Pass) GeneratePasswordHash() error {
	if p.Password == "" {
		return er.ErrPasswordIsEmpty
	}

	hash, err := generatePasswordHash(p.Password)

	if err != nil {
		return err
	}

	p.PasswordHash = hash
	return nil
}

// CheckPassword ...
func (p *Pass) CheckPassword(password string) error {
	if p.PasswordHash == "" {
		err := p.GeneratePasswordHash()
		if err != nil {
			return err
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(password))

	if err != nil {
		return err
	}

	return nil
}

func generatePasswordHash(pass string) (string, error) {
	if pass == "" {
		return "", er.ErrPasswordIsEmpty
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}
