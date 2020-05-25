package db

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/models"
)

// LoggedDB represents logged users
// it's ok to be in-memory pool
// don't forget use Locker
type LoggedDB interface {
	sync.Locker

	// Login must return valid *User if (error == ErrAlreadyLogged)
	// with other errors must return (nil, err)
	// also should generate valid LoginToken and LastActivity
	Login(username string) (*models.LoggedUser, error)

	Get(username string) (*models.LoggedUser, error)

	Logout(username string) error
}
