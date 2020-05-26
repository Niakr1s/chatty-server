package chat

import (
	"time"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
)

// NotifyChat ...
type NotifyChat struct {
	db.Chat

	notifyCh chan<- events.Event
}

// NewNotifyChat ...
func NewNotifyChat(chat db.Chat, ch chan<- events.Event) *NotifyChat {
	return &NotifyChat{chat, ch}
}

// AddUser ...
func (c *NotifyChat) AddUser(username string) error {
	err := c.Chat.AddUser(username)

	if err != nil {
		return err
	}

	c.notifyUserJoined(username, c.Chat.ChatName(), time.Now())

	return nil
}

// RemoveUser ...
func (c *NotifyChat) RemoveUser(username string) error {
	err := c.Chat.RemoveUser(username)

	if err != nil {
		return err
	}

	c.notifyUserLeaved(username, c.Chat.ChatName(), time.Now())

	return nil
}

func (c *NotifyChat) notifyUserJoined(username, chatname string, t time.Time) {
	go func() {
		c.notifyCh <- events.NewChatJoinEvent(username, chatname, t)
	}()
}

func (c *NotifyChat) notifyUserLeaved(username, chatname string, t time.Time) {
	go func() {
		c.notifyCh <- events.NewChatLeaveEvent(username, chatname, t)
	}()
}
