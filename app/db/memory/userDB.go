package memory

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/models"
)

// UserDB ...
type UserDB struct {
	sync.RWMutex

	users map[string]models.FullUser
}

// NewUserDB ...
func NewUserDB() *UserDB {
	return &UserDB{users: make(map[string]models.FullUser)}
}

// Store ...
func (d *UserDB) Store(u models.FullUser) error {
	d.Lock()
	defer d.Unlock()

	if _, ok := d.users[u.UserName]; ok {
		return er.ErrUserAlreadyRegistered
	}

	d.users[u.UserName] = u
	return nil
}

// Update ...
func (d *UserDB) Update(u models.FullUser) error {
	d.Lock()
	defer d.Unlock()

	d.users[u.UserName] = u
	return nil
}

// Get ...
func (d *UserDB) Get(username string) (models.FullUser, error) {
	d.RLock()
	defer d.RUnlock()

	u, ok := d.users[username]
	if !ok {
		return models.FullUser{}, er.ErrUserNotFound
	}

	return u, nil
}
