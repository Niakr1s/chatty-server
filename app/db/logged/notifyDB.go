package logged

import (
	"time"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/niakr1s/chatty-server/app/pool/events"
)

// NotifyDB ...
type NotifyDB struct {
	db.LoggedDB

	notifyCh chan<- events.Event
}

// NewNotifyDB ...
func NewNotifyDB(db db.LoggedDB, ch chan<- events.Event) *NotifyDB {
	return &NotifyDB{db, ch}
}

// Login ...
func (d *NotifyDB) Login(username string) (*models.LoggedUser, error) {
	u, err := d.LoggedDB.Login(username)

	if err != nil {
		return u, err
	}

	d.notifyLogin(username, u.LastActivity)

	return u, err
}

// Logout ...
func (d *NotifyDB) Logout(username string) error {
	err := d.LoggedDB.Logout(username)

	if err != nil {
		return err
	}

	d.notifyLogout(username, time.Now())

	return err
}

func (d *NotifyDB) notifyLogin(username string, t time.Time) {
	go func() {
		if d.notifyCh != nil {
			d.notifyCh <- events.NewLoginEvent(username, "", t)
		}
	}()
}

func (d *NotifyDB) notifyLogout(username string, t time.Time) {
	go func() {
		if d.notifyCh != nil {
			d.notifyCh <- events.NewLogoutEvent(username, "", t)
		}
	}()
}
