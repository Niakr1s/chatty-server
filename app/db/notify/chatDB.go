package notify

import (
	"time"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
)

// ChatDB is wrapper above ChatDB
type ChatDB struct {
	db.ChatDB

	notifyCh chan<- events.Event
}

// NewChatDB ...
func NewChatDB(db db.ChatDB, ch chan<- events.Event) *ChatDB {
	return &ChatDB{ChatDB: db, notifyCh: ch}
}

// StartListeningToEvents starts to listen to events from input channel.
// Don't call it twice!
func (d *ChatDB) StartListeningToEvents(ch <-chan events.Event) {
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
func (d *ChatDB) Add(chatname string) (db.Chat, error) {
	c, err := d.ChatDB.Add(chatname)

	if err != nil {
		return c, err
	}

	// converting chat to notifyChat
	c = NewChat(c, d.notifyCh)

	d.notifyChatCreated(chatname, time.Now().UTC())

	return c, nil
}

// Get ...
func (d *ChatDB) Get(chatname string) (db.Chat, error) {
	c, err := d.ChatDB.Get(chatname)
	if err != nil {
		return c, err
	}
	return NewChat(c, d.notifyCh), nil
}

// Remove ...
func (d *ChatDB) Remove(chatname string) error {
	err := d.ChatDB.Remove(chatname)

	if err != nil {
		return err
	}

	d.notifyChatRemoved(chatname, time.Now().UTC())

	return nil
}

// GetChats ...
func (d *ChatDB) GetChats() []db.Chat {
	chats := d.ChatDB.GetChats()
	for i, chat := range chats {
		chats[i] = NewChat(chat, d.notifyCh)
	}
	return chats
}

func (d *ChatDB) notifyChatCreated(chatname string, t time.Time) {
	go func() {
		d.notifyCh <- events.NewChatCreatedEvent(chatname, t)
	}()
}

func (d *ChatDB) notifyChatRemoved(chatname string, t time.Time) {
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
