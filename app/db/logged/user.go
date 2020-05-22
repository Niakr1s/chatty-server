package logged

import (
	"server2/app/random"
	"time"
)

// User ...
type User struct {
	Name string

	lastActivity time.Time
	loginToken   string
}

// NewUser creates new user, updates last activity and updates login token
func NewUser(username string) *User {
	res := &User{Name: username}
	res.UpdateLastActivity(time.Now())
	res.UpdateLoginToken()
	return res
}

// UpdateLastActivity ...
func (u *User) UpdateLastActivity(t time.Time) *User {
	u.lastActivity = t
	return u
}

// UpdateLoginToken ...
func (u *User) UpdateLoginToken() *User {
	u.loginToken = random.StrGen.RandomStr()
	return u
}
