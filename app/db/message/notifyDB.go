package message

import (
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/niakr1s/chatty-server/app/pool/events"
)

// NotifyDB ...
type NotifyDB struct {
	db.MessageDB

	notifyCh chan<- events.Event
}

// NewNotifyDB ...
func NewNotifyDB(db db.MessageDB, ch chan<- events.Event) *NotifyDB {
	return &NotifyDB{db, ch}
}

// Post ...
func (d *NotifyDB) Post(msg *models.Message) error {
	err := d.MessageDB.Post(msg)

	if err != nil {
		return err
	}

	d.notifyNewMessage(msg)

	return nil
}

func (d *NotifyDB) notifyNewMessage(msg *models.Message) {
	go func() {
		d.notifyCh <- events.NewMessageEvent(msg)
	}()
}
