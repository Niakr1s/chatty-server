package models

// User ...
type User struct {
	UserName string `json:"user" validate:"required"`
}

// NewUser ...
func NewUser(username string) User {
	return User{UserName: username}
}
