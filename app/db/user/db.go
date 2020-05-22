package user

import (
	"errors"
	"server2/app/models"
)

// DB should be persistend storage
type DB interface {
	// should update user's ID
	Store(u *models.User) error

	Get(id uint) (models.User, error)
}

// Errors
var (
	ErrUserNotExist = errors.New("user doesn't exist in user db")
)
