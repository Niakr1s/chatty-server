package db

import (
	"github.com/niakr1s/chatty-server/app/models"
)

// UserDB should be persistend storage
type UserDB interface {
	Store(u models.FullUser) error
	Update(u models.FullUser) error

	Get(username string) (models.FullUser, error)
}

// IsUserVerified gets verified status.
func IsUserVerified(db UserDB, username string) bool {
	u, err := db.Get(username)
	if err != nil {
		return false
	}
	return u.Verified
}
