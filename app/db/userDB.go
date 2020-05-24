package db

import (
	"github.com/niakr1s/chatty-server/app/models"
)

// UserDB should be persistend storage
type UserDB interface {
	Store(u *models.FullUser) error

	Get(username string) (models.FullUser, error)
}
