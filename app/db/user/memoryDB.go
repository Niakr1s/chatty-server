package user

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/models"
)

// MemoryDB ...
type MemoryDB struct {
	sync.RWMutex

	users map[string]*models.FullUser
}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{users: make(map[string]*models.FullUser)}
}

// Store ...
func (d *MemoryDB) Store(u *models.FullUser) error {
	d.Lock()
	defer d.Unlock()

	if _, ok := d.users[u.Name]; ok {
		return er.ErrUserAlreadyRegistered
	}

	d.users[u.Name] = u
	return nil
}

// Get ...
func (d *MemoryDB) Get(username string) (models.FullUser, error) {
	d.RLock()
	defer d.RUnlock()

	u, ok := d.users[username]
	if !ok {
		return models.FullUser{}, er.ErrUserNotFound
	}

	return *u, nil
}
