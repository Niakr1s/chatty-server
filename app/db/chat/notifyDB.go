package chat

import (
	"time"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/pool/events"
)

// NotifyDB is wrapper above ChatDB
type NotifyDB struct {
	db.ChatDB

	notifyCh chan<- events.Event
}

// NewNotifyDB ...
func NewNotifyDB(db db.ChatDB, ch chan<- events.Event) *NotifyDB {
	return &NotifyDB{ChatDB: db, notifyCh: ch}
}

// Add ...
func (d *NotifyDB) Add(chatname string) (db.Chat, error) {
	c, err := d.ChatDB.Add(chatname)

	if err != nil {
		return c, err
	}

	// converting chat to notifyChat
	c = NewNotifyChat(c, d.notifyCh)

	d.notifyChatCreated(chatname, time.Now())

	return c, nil
}

// Remove ...
func (d *NotifyDB) Remove(chatname string) error {
	err := d.ChatDB.Remove(chatname)

	if err != nil {
		return err
	}

	d.notifyChatRemoved(chatname, time.Now())

	return nil
}

func (d *NotifyDB) notifyChatCreated(chatname string, t time.Time) {
	go func() {
		d.notifyCh <- events.NewChatCreatedEvent(chatname, t)
	}()
}

func (d *NotifyDB) notifyChatRemoved(chatname string, t time.Time) {
	go func() {
		d.notifyCh <- events.NewChatRemovedEvent(chatname, t)
	}()
}
