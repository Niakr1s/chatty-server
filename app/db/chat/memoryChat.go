package chat

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
)

// MemoryChat ...
type MemoryChat struct {
	sync.Mutex

	Name  string
	users map[string]struct{}

	notifyCh chan<- events.Event
}

// NewMemoryChat ...
func NewMemoryChat(chatname string) *MemoryChat {
	return &MemoryChat{Name: chatname, users: make(map[string]struct{})}
}

// WithNotifyCh ...
func (c *MemoryChat) WithNotifyCh(ch chan<- events.Event) *NotifyChat {
	return NewNotifyChat(c, ch)
}

// ChatName ...
func (c *MemoryChat) ChatName() string {
	return c.Name
}

// AddUser ...
func (c *MemoryChat) AddUser(username string) error {
	if _, ok := c.users[username]; ok {
		return er.ErrAlreadyInChat
	}

	c.users[username] = struct{}{}
	return nil
}

// RemoveUser ...
func (c *MemoryChat) RemoveUser(username string) error {
	if _, ok := c.users[username]; !ok {
		return er.ErrNotInChat
	}
	delete(c.users, username)

	return nil
}

// IsInChat ...
func (c *MemoryChat) IsInChat(username string) bool {
	_, ok := c.users[username]
	return ok
}

// GetUsers ...
func (c *MemoryChat) GetUsers() []models.User {
	res := make([]models.User, 0, len(c.users))
	for u := range c.users {
		res = append(res, models.User{UserName: u})
	}
	return res
}
