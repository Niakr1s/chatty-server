package memory

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/models"
)

// Chat ...
type Chat struct {
	sync.Mutex

	Name  string
	users map[string]struct{}
}

// NewChat ...
func NewChat(chatname string) *Chat {
	return &Chat{Name: chatname, users: make(map[string]struct{})}
}

// ChatName ...
func (c *Chat) ChatName() string {
	return c.Name
}

// AddUser ...
func (c *Chat) AddUser(username string) error {
	if _, ok := c.users[username]; ok {
		return er.ErrAlreadyInChat
	}

	c.users[username] = struct{}{}
	return nil
}

// RemoveUser ...
func (c *Chat) RemoveUser(username string) error {
	if _, ok := c.users[username]; !ok {
		return er.ErrNotInChat
	}
	delete(c.users, username)

	return nil
}

// IsInChat ...
func (c *Chat) IsInChat(username string) bool {
	_, ok := c.users[username]
	return ok
}

// GetUsers ...
func (c *Chat) GetUsers() []models.User {
	res := make([]models.User, 0, len(c.users))
	for u := range c.users {
		res = append(res, models.User{UserName: u})
	}
	return res
}
