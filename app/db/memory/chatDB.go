package memory

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/er"
)

// ChatDB ...
type ChatDB struct {
	sync.Mutex

	chats map[string]db.Chat
}

// NewChatDB ...
func NewChatDB() *ChatDB {
	return &ChatDB{chats: make(map[string]db.Chat)}
}

// Add ...
func (d *ChatDB) Add(chatname string) (db.Chat, error) {
	if c, ok := d.chats[chatname]; ok {
		return c, er.ErrChatAlreadyExists
	}
	c := NewChat(chatname)
	d.chats[chatname] = c

	return c, nil
}

// Get ...
func (d *ChatDB) Get(chatname string) (db.Chat, error) {
	if c, ok := d.chats[chatname]; ok {
		return c, nil
	}
	return nil, er.ErrNoSuchChat
}

// Remove ...
func (d *ChatDB) Remove(chatname string) error {
	if _, ok := d.chats[chatname]; !ok {
		return er.ErrNoSuchChat
	}
	delete(d.chats, chatname)

	return nil
}

// GetChats ...
func (d *ChatDB) GetChats() []db.Chat {
	res := make([]db.Chat, 0, len(d.chats))

	for _, c := range d.chats {
		res = append(res, c)
	}

	return res
}
