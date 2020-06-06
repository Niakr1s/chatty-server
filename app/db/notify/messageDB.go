package notify

import (
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
)

// MessageDB ...
type MessageDB struct {
	db.MessageDB

	notifyCh chan<- events.Event
}

// NewMessageDB ...
func NewMessageDB(db db.MessageDB, ch chan<- events.Event) *MessageDB {
	return &MessageDB{db, ch}
}

// Post ...
func (d *MessageDB) Post(msg *models.Message) error {
	err := d.MessageDB.Post(msg)

	if err != nil {
		return err
	}

	d.notifyNewMessage(msg)

	return nil
}

func (d *MessageDB) notifyNewMessage(msg *models.Message) {
	go func() {
		d.notifyCh <- events.NewMessageEvent(msg)
	}()
}
