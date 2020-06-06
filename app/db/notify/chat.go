package notify

import (
	"time"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
)

// Chat ...
type Chat struct {
	db.Chat
	logged db.LoggedDB

	notifyCh chan<- events.Event
}

// NewChat ...
func NewChat(chat db.Chat, logged db.LoggedDB, ch chan<- events.Event) *Chat {
	return &Chat{chat, logged, ch}
}

// AddUser ...
func (c *Chat) AddUser(username string) error {
	err := c.Chat.AddUser(username)

	if err != nil {
		return err
	}

	c.notifyUserJoined(username, c.Chat.ChatName(), time.Now().UTC())

	return nil
}

// RemoveUser ...
func (c *Chat) RemoveUser(username string) error {
	err := c.Chat.RemoveUser(username)

	if err != nil {
		return err
	}

	c.notifyUserLeaved(username, c.Chat.ChatName(), time.Now().UTC())

	return nil
}

func (c *Chat) notifyUserJoined(username, chatname string, t time.Time) {
	go func() {
		c.notifyCh <- events.NewChatJoinEvent(username, chatname, t)
	}()
}

func (c *Chat) notifyUserLeaved(username, chatname string, t time.Time) {
	go func() {
		c.notifyCh <- events.NewChatLeaveEvent(username, chatname, t)
	}()
}
