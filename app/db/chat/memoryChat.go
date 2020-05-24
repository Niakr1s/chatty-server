package chat

import (
	"time"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/pool/events"
)

// MemoryChat ...
type MemoryChat struct {
	Name  string
	users map[string]struct{}

	notifyCh chan<- events.Event
}

// NewMemoryChat ...
func NewMemoryChat(chatname string) *MemoryChat {
	return &MemoryChat{Name: chatname, users: make(map[string]struct{})}
}

// WithNotifyCh ...
func (c *MemoryChat) WithNotifyCh(ch chan<- events.Event) *MemoryChat {
	c.notifyCh = ch
	return c
}

// AddUser ...
func (c *MemoryChat) AddUser(username string) error {
	if _, ok := c.users[username]; ok {
		return er.ErrAlreadyInChat
	}

	c.notifyUserJoined(username, c.Name, time.Now())

	c.users[username] = struct{}{}
	return nil
}

// RemoveUser ...
func (c *MemoryChat) RemoveUser(username string) error {
	if _, ok := c.users[username]; !ok {
		return er.ErrNotInChat
	}
	delete(c.users, username)

	c.notifyUserLeaved(username, c.Name, time.Now())

	return nil
}

// IsInChat ...
func (c *MemoryChat) IsInChat(username string) bool {
	_, ok := c.users[username]
	return ok
}

func (c *MemoryChat) notifyUserJoined(username, chatname string, t time.Time) {
	go func() {
		if c.notifyCh != nil {
			c.notifyCh <- events.NewChatJoinEvent(username, chatname, t)
		}
	}()
}

func (c *MemoryChat) notifyUserLeaved(username, chatname string, t time.Time) {
	go func() {
		if c.notifyCh != nil {
			c.notifyCh <- events.NewChatLeaveEvent(username, chatname, t)
		}
	}()
}
