package notify

import (
	"time"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
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
		c.logged.Lock()
		status := models.UserStatus{}
		if loggedU, err := c.logged.Get(username); err == nil {
			status = loggedU.UserStatus
		}
		c.logged.Unlock()
		c.notifyCh <- events.NewChatJoinEvent(username, chatname, t).WithStatus(status)
	}()
}

func (c *Chat) notifyUserLeaved(username, chatname string, t time.Time) {
	go func() {
		c.notifyCh <- events.NewChatLeaveEvent(username, chatname, t)
	}()
}
