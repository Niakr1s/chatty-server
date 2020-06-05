package chat

import (
	"time"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
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

// StartListeningToEvents starts to listen to events from input channel.
// Don't call it twice!
func (d *NotifyDB) StartListeningToEvents(ch <-chan events.Event) {
	go func() {
		for {
			event := <-ch
			switch e := event.(type) {
			case *events.LogoutEvent:
				logoutUserFromAllChats(e.UserName, d)
			}
		}
	}()
}

// Add ...
func (d *NotifyDB) Add(chatname string) (db.Chat, error) {
	c, err := d.ChatDB.Add(chatname)

	if err != nil {
		return c, err
	}

	// converting chat to notifyChat
	c = NewNotifyChat(c, d.notifyCh)

	d.notifyChatCreated(chatname, time.Now().UTC())

	return c, nil
}

// Get ...
func (d *NotifyDB) Get(chatname string) (db.Chat, error) {
	c, err := d.ChatDB.Get(chatname)
	if err != nil {
		return c, err
	}
	return NewNotifyChat(c, d.notifyCh), nil
}

// Remove ...
func (d *NotifyDB) Remove(chatname string) error {
	err := d.ChatDB.Remove(chatname)

	if err != nil {
		return err
	}

	d.notifyChatRemoved(chatname, time.Now().UTC())

	return nil
}

// GetChats ...
func (d *NotifyDB) GetChats() []db.Chat {
	chats := d.ChatDB.GetChats()
	for i, chat := range chats {
		chats[i] = NewNotifyChat(chat, d.notifyCh)
	}
	return chats
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

// logoutUserFromAllChats forcefully logouts user from all chats.
// It uses locks, beware of it.
// It is used in NotifyDB to logout users on LogoutEvent
func logoutUserFromAllChats(username string, chatDB db.ChatDB) {
	chatDB.Lock()
	defer chatDB.Unlock()

	chats := chatDB.GetChats()
	for _, chat := range chats {
		chat.Lock()
		chat.RemoveUser(username)
		chat.Unlock()
	}
}
