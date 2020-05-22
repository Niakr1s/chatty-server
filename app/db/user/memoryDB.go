package user

import (
	"server2/app/models"
	"sync"
)

// MemoryDB ...
type MemoryDB struct {
	sync.RWMutex

	counter uint
	users   map[uint]*models.User
}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{users: make(map[uint]*models.User)}
}

// Store ...
func (d *MemoryDB) Store(u *models.User) error {
	d.Lock()
	defer d.Unlock()

	d.users[d.counter] = u
	u.ID = d.counter
	d.counter++

	return nil
}

// Get ...
func (d *MemoryDB) Get(id uint) (models.User, error) {
	d.RLock()
	defer d.RUnlock()

	u, ok := d.users[id]
	if !ok {
		return models.User{}, ErrUserNotExist
	}

	return *u, nil
}
