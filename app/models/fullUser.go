package models

import (
	"github.com/niakr1s/chatty-server/app/internal/validator"
)

// FullUser ...
type FullUser struct {
	User
	UserStatus
	Email
	Pass
}

// NewFullUser ...
func NewFullUser(username, email, password string) FullUser {
	return FullUser{User: User{UserName: username}, Email: Email{Address: email}, Pass: Pass{Password: password}}
}

// ValidateBeforeStoring used before storing in database
func (u *FullUser) ValidateBeforeStoring() error {
	return validator.Validate.Struct(*u)
}
