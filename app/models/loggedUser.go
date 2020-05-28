package models

import (
	"time"

	"github.com/niakr1s/chatty-server/app/internal/random"
)

// LoggedUser ...
type LoggedUser struct {
	User

	LastActivity time.Time `validate:"required"`
	LoginToken   string    `validate:"required"`
}

// NewLoggedUser creates new user, updates last activity and updates login token
func NewLoggedUser(username string) *LoggedUser {
	res := &LoggedUser{User: User{Name: username}}
	res.UpdateLastActivity(time.Now())
	res.UpdateLoginToken()
	return res
}

// UpdateLastActivity ...
func (u *LoggedUser) UpdateLastActivity(t time.Time) *LoggedUser {
	u.LastActivity = t
	return u
}

// UpdateLoginToken ...
func (u *LoggedUser) UpdateLoginToken() *LoggedUser {
	u.LoginToken = random.StrGen.RandomStr()
	return u
}
