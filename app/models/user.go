package models

// User ...
type User struct {
	Name string `json:"name" validate:"required"`
}
