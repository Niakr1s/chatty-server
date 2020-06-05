package db

import (
	"sync"
	"time"

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

	GetLoggedUsers() []string
}

// StartCleanInactiveUsers ...
func StartCleanInactiveUsers(d LoggedDB, each time.Duration, inactivityTimeout time.Duration) {
	go func() {
		for {
			<-time.After(each)
			cleanInactiveUsers(d, inactivityTimeout)
		}
	}()
}

// CleanInactiveUsers ...
func cleanInactiveUsers(d LoggedDB, inactivityTimeout time.Duration) {
	d.Lock()
	defer d.Unlock()

	users := d.GetLoggedUsers()
	now := time.Now().UTC()
	for _, username := range users {
		user, err := d.Get(username)
		if err != nil {
			continue
		}
		if diff := now.Sub(user.LastActivity); diff > inactivityTimeout {
			d.Logout(username)
		}
	}
}
