package chat

import "server2/app/er"

// Chat ...
type Chat struct {
	Name  string
	users map[string]struct{}
}

// NewChat ...
func NewChat(chatname string) *Chat {
	return &Chat{Name: chatname, users: make(map[string]struct{})}
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
