package notify

import (
	"time"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
)

// LoggedDB ...
type LoggedDB struct {
	db.LoggedDB

	notifyCh chan<- events.Event
}

// NewLoggedDB ...
func NewLoggedDB(db db.LoggedDB, ch chan<- events.Event) *LoggedDB {
	return &LoggedDB{db, ch}
}

// Login ...
func (d *LoggedDB) Login(username string) (*models.LoggedUser, error) {
	u, err := d.LoggedDB.Login(username)

	if err != nil {
		return u, err
	}

	d.notifyLogin(username, u.LastActivity)

	return u, nil
}

// Logout ...
func (d *LoggedDB) Logout(username string) error {
	err := d.LoggedDB.Logout(username)

	if err != nil {
		return err
	}

	d.notifyLogout(username, time.Now().UTC())

	return nil
}

func (d *LoggedDB) notifyLogin(username string, t time.Time) {
	go func() {
		d.notifyCh <- events.NewLoginEvent(username, "", t)
	}()
}

func (d *LoggedDB) notifyLogout(username string, t time.Time) {
	go func() {
		d.notifyCh <- events.NewLogoutEvent(username, "", t)
	}()
}
