package memory

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/er"
)

type users map[string]struct{}

// ChatDB ...
type ChatDB struct {
	sync.Mutex

	chats map[string]users
}

// NewChatDB ...
func NewChatDB() *ChatDB {
	return &ChatDB{chats: make(map[string]users)}
}

// Add ...
func (d *ChatDB) Add(chatname string) error {
	d.Lock()
	defer d.Unlock()

	return d.add(chatname)
}

// Remove ...
func (d *ChatDB) Remove(chatname string) error {
	d.Lock()
	defer d.Unlock()

	return d.remove(chatname)
}

// GetChats ...
func (d *ChatDB) GetChats() []string {
	d.Lock()
	defer d.Unlock()

	return d.getChats()
}

// AddUser ...
func (d *ChatDB) AddUser(chatname, username string) error {
	d.Lock()
	defer d.Unlock()

	return d.addUser(chatname, username)
}

// RemoveUser ...
func (d *ChatDB) RemoveUser(chatname, username string) error {
	d.Lock()
	defer d.Unlock()

	return d.removeUser(chatname, username)
}

// IsInChat ...
func (d *ChatDB) IsInChat(chatname, username string) bool {
	d.Lock()
	defer d.Unlock()

	return d.isInChat(chatname, username)
}

// GetUsers ...
func (d *ChatDB) GetUsers(chatname string) []string {
	d.Lock()
	defer d.Unlock()

	return d.getUsers(chatname)
}

// concurrent-unsafe methods

// Add ...
func (d *ChatDB) add(chatname string) error {
	if _, ok := d.chats[chatname]; ok {
		return er.ErrChatAlreadyExists
	}
	d.chats[chatname] = make(users)

	return nil
}

// Remove ...
func (d *ChatDB) remove(chatname string) error {
	if _, ok := d.chats[chatname]; !ok {
		return er.ErrNoSuchChat
	}
	delete(d.chats, chatname)

	return nil
}

// GetChats ...
func (d *ChatDB) getChats() []string {
	res := make([]string, 0, len(d.chats))

	for c := range d.chats {
		res = append(res, c)
	}

	return res
}

func (d *ChatDB) addUser(chatname, username string) error {
	c, ok := d.chats[chatname]
	if !ok {
		return er.ErrNoSuchChat
	}

	if _, ok := c[username]; ok {
		return er.ErrAlreadyInChat
	}

	d.chats[chatname][username] = struct{}{}
	return nil
}

func (d *ChatDB) removeUser(chatname, username string) error {
	c, ok := d.chats[chatname]
	if !ok {
		return er.ErrNoSuchChat
	}

	if _, ok := c[username]; !ok {
		return er.ErrNotInChat
	}

	delete(d.chats[chatname], username)
	return nil
}

func (d *ChatDB) isInChat(chatname, username string) bool {
	c, ok := d.chats[chatname]
	if !ok {
		return false
	}

	if _, ok := c[username]; !ok {
		return false
	}

	return true
}

func (d *ChatDB) getUsers(chatname string) []string {
	c, ok := d.chats[chatname]
	if !ok {
		return []string{}
	}

	res := make([]string, len(c))
	for u := range c {
		res = append(res, u)
	}
	return res
}
