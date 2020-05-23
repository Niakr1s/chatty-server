package logged

import (
	"sync"
	"time"
)

// MemoryDB ...
type MemoryDB struct {
	sync.Mutex

	users map[string]*User
}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{users: make(map[string]*User)}
}

// StartCleanInactiveUsers ...
func (d *MemoryDB) StartCleanInactiveUsers(each time.Duration, inactivityTimeout time.Duration) {
	go func() {
		for {
			<-time.After(each)
			d.cleanInactiveUsers(inactivityTimeout)
		}
	}()
}

func (d *MemoryDB) cleanInactiveUsers(inactivityTimeout time.Duration) {
	d.Lock()
	defer d.Unlock()

	now := time.Now()

	for n, u := range d.users {
		if diff := now.Sub(u.lastActivity); diff > inactivityTimeout {
			d.Logout(n)
		}
	}
}

// Login ...
func (d *MemoryDB) Login(username string) (*User, error) {
	u, ok := d.users[username]

	if ok {
		return u, ErrAlreadyLogged
	}

	u = NewUser(username)
	d.users[username] = u
	return u, nil
}

// Get ...
func (d *MemoryDB) Get(username string) (*User, error) {
	if u, ok := d.users[username]; ok {
		return u, nil
	}

	return nil, ErrNotLogged
}

// Logout ...
func (d *MemoryDB) Logout(username string) error {
	if _, ok := d.users[username]; !ok {
		return ErrNotLogged
	}
	delete(d.users, username)
	return nil
}
