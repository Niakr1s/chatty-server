package logged

import "sync"

// DB represents logged users
// it's ok to be in-memory pool
// don't forget use Locker
type DB interface {
	sync.Locker

	// Login must return valid *User if (error == ErrAlreadyLogged)
	// with other errors must return (nil, err)
	Login(username string) (*User, error)

	Get(username string) (*User, error)

	Logout(username string) error
}
