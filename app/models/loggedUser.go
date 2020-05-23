package models

import (
	"server2/app/random"
	"time"
)

// LoggedUser ...
type LoggedUser struct {
	Name string

	LastActivity time.Time
	LoginToken   string
}

// NewLoggedUser creates new user, updates last activity and updates login token
func NewLoggedUser(username string) *LoggedUser {
	res := &LoggedUser{Name: username}
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