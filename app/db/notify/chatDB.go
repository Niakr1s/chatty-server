package notify

import (
	"time"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
)

// ChatDB is wrapper above ChatDB
type ChatDB struct {
	db.ChatDB
	logged db.LoggedDB

	notifyCh chan<- events.Event
}

// NewChatDB ...
func NewChatDB(db db.ChatDB, logged db.LoggedDB, ch chan<- events.Event) *ChatDB {
	return &ChatDB{db, logged, ch}
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
func (d *ChatDB) Add(chatname string) error {
	err := d.ChatDB.Add(chatname)
	if err != nil {
		return err
	}

	d.notifyChatCreated(chatname, time.Now().UTC())

	return nil
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

// AddUser ...
func (d *ChatDB) AddUser(chatname, username string) error {
	err := d.ChatDB.AddUser(chatname, username)
	if err != nil {
		return err
	}

	d.notifyUserJoined(chatname, username, time.Now().UTC())
	return err
}

// RemoveUser ...
func (d *ChatDB) RemoveUser(chatname, username string) error {
	err := d.ChatDB.RemoveUser(chatname, username)
	if err != nil {
		return err
	}
	d.notifyUserLeaved(chatname, username, time.Now().UTC())
	return err
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

func (d *ChatDB) notifyUserJoined(chatname, username string, t time.Time) {
	go func() {
		status := models.UserStatus{}
		if loggedU, err := d.logged.Get(username); err == nil {
			status = loggedU.UserStatus
		}
		d.notifyCh <- events.NewChatJoinEvent(username, chatname, t).WithStatus(status)
	}()
}

func (d *ChatDB) notifyUserLeaved(chatname, username string, t time.Time) {
	go func() {
		d.notifyCh <- events.NewChatLeaveEvent(username, chatname, t)
	}()
}

// logoutUserFromAllChats forcefully logouts user from all chats.
// It uses locks, beware of it.
// It is used in NotifyDB to logout users on LogoutEvent
func logoutUserFromAllChats(username string, chatDB db.ChatDB) {
	chats := chatDB.GetChats()
	for _, chat := range chats {
		chatDB.RemoveUser(chat, username)
	}
}
