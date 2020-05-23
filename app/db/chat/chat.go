package chat

import (
	"server2/app/er"
	"server2/app/pool/events"
	"time"
)

// Chat ...
type Chat struct {
	Name  string
	users map[string]struct{}

	notifyCh chan<- events.Event
}

// NewChat ...
func NewChat(chatname string) *Chat {
	return &Chat{Name: chatname, users: make(map[string]struct{})}
}

// WithNotifyCh ...
func (c *Chat) WithNotifyCh(ch chan<- events.Event) *Chat {
	c.notifyCh = ch
	return c
}

// AddUser ...
func (c *Chat) AddUser(username string) error {
	if _, ok := c.users[username]; ok {
		return er.ErrAlreadyInChat
	}

	c.notifyUserJoined(username, c.Name, time.Now())

	c.users[username] = struct{}{}
	return nil
}

// RemoveUser ...
func (c *Chat) RemoveUser(username string) error {
	if _, ok := c.users[username]; !ok {
		return er.ErrNotInChat
	}
	delete(c.users, username)

	c.notifyUserLeaved(username, c.Name, time.Now())

	return nil
}

// IsInChat ...
func (c *Chat) IsInChat(username string) bool {
	_, ok := c.users[username]
	return ok
}

func (c *Chat) notifyUserJoined(username, chatname string, t time.Time) {
	go func() {
		if c.notifyCh != nil {
			c.notifyCh <- events.NewChatJoinEvent(username, chatname, t)
		}
	}()
}

func (c *Chat) notifyUserLeaved(username, chatname string, t time.Time) {
	go func() {
		if c.notifyCh != nil {
			c.notifyCh <- events.NewChatLeaveEvent(username, chatname, t)
		}
	}()
}
