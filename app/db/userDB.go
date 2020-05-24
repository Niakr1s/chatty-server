package db

import (
	"github.com/niakr1s/chatty-server/app/models"
)

// UserDB should be persistend storage
type UserDB interface {
	// should update user's ID
	Store(u *models.User) error

	Get(id uint) (models.User, error)
}
