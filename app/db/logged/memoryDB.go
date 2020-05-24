package logged

import (
	"sync"
	"time"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/niakr1s/chatty-server/app/pool/events"
)

// MemoryDB ...
type MemoryDB struct {
	sync.Mutex

	users map[string]*models.LoggedUser

	notifyCh chan<- events.Event
}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{users: make(map[string]*models.LoggedUser)}
}

// WithNotifyCh ...
func (d *MemoryDB) WithNotifyCh(ch chan<- events.Event) *MemoryDB {
	d.notifyCh = ch
	return d
}

// StartCleanInactiveUsers ...
func (d *MemoryDB) StartCleanInactiveUsers(each time.Duration, inactivityTimeout time.Duration) {
	go func() {
		for {
			<-time.After(each)
			d.Lock()
			d.cleanInactiveUsers(inactivityTimeout)
			d.Unlock()
		}
	}()
}

func (d *MemoryDB) cleanInactiveUsers(inactivityTimeout time.Duration) {
	now := time.Now()

	for n, u := range d.users {
		if diff := now.Sub(u.LastActivity); diff > inactivityTimeout {
			d.Logout(n)
		}
	}
}

// Login ...
func (d *MemoryDB) Login(username string) (*models.LoggedUser, error) {
	u, ok := d.users[username]

	if ok {
		return u, er.ErrAlreadyLogged
	}

	u = models.NewLoggedUser(username)
	d.users[username] = u

	d.notifyLogin(username, u.LastActivity)

	return u, nil
}

// Get ...
func (d *MemoryDB) Get(username string) (*models.LoggedUser, error) {
	if u, ok := d.users[username]; ok {
		return u, nil
	}

	return nil, er.ErrNotLogged
}

// Logout ...
func (d *MemoryDB) Logout(username string) error {
	if _, ok := d.users[username]; !ok {
		return er.ErrNotLogged
	}
	delete(d.users, username)

	d.notifyLogout(username, time.Now())

	return nil
}

func (d *MemoryDB) notifyLogin(username string, t time.Time) {
	go func() {
		if d.notifyCh != nil {
			d.notifyCh <- events.NewLoginEvent(username, "", t)
		}
	}()
}

func (d *MemoryDB) notifyLogout(username string, t time.Time) {
	go func() {
		if d.notifyCh != nil {
			d.notifyCh <- events.NewLogoutEvent(username, "", t)
		}
	}()
}
