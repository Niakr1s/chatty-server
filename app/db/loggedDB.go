package db

import (
	"server2/app/models"
	"sync"
)

// LoggedDB represents logged users
// it's ok to be in-memory pool
// don't forget use Locker
type LoggedDB interface {
	sync.Locker

	// Login must return valid *User if (error == ErrAlreadyLogged)
	// with other errors must return (nil, err)
	Login(username string) (*models.LoggedUser, error)

	Get(username string) (*models.LoggedUser, error)

	Logout(username string) error
}
