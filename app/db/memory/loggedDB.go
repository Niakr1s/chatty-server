package memory

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/models"
)

// LoggedDB ...
type LoggedDB struct {
	sync.Mutex

	users map[string]*models.LoggedUser
}

// NewLoggedDB ...
func NewLoggedDB() *LoggedDB {
	return &LoggedDB{users: make(map[string]*models.LoggedUser)}
}

// GetLoggedUsers ...
func (d *LoggedDB) GetLoggedUsers() []string {
	res := make([]string, 0, len(d.users))
	for u := range d.users {
		res = append(res, u)
	}
	return res
}

// Login ...
func (d *LoggedDB) Login(username string) (*models.LoggedUser, error) {
	u, ok := d.users[username]

	if ok {
		return u, er.ErrAlreadyLogged
	}

	u = models.NewLoggedUser(username)
	d.users[username] = u

	return u, nil
}

// Update ...
func (d *LoggedDB) Update(user *models.LoggedUser) error {
	if _, ok := d.users[user.UserName]; !ok {
		return er.ErrNotLogged
	}
	d.users[user.UserName] = user
	return nil
}

// Get ...
func (d *LoggedDB) Get(username string) (*models.LoggedUser, error) {
	if u, ok := d.users[username]; ok {
		return u, nil
	}

	return nil, er.ErrNotLogged
}

// Logout ...
func (d *LoggedDB) Logout(username string) error {
	if _, ok := d.users[username]; !ok {
		return er.ErrNotLogged
	}
	delete(d.users, username)

	return nil
}
