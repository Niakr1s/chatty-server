package chat

import (
	"server2/app/er"
	"server2/app/models"
	"sync"
)

// MemoryDB ...
type MemoryDB struct {
	sync.Mutex

	chats map[string]*models.Chat
}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{chats: make(map[string]*models.Chat)}
}

// Add ...
func (d *MemoryDB) Add(chatname string) (*models.Chat, error) {
	if c, ok := d.chats[chatname]; ok {
		return c, er.ErrChatAlreadyExists
	}
	c := models.NewChat(chatname)
	d.chats[chatname] = c
	return c, nil
}

// Get ...
func (d *MemoryDB) Get(chatname string) (*models.Chat, error) {
	if c, ok := d.chats[chatname]; ok {
		return c, nil
	}
	return nil, er.ErrNoSuchChat
}

// Remove ...
func (d *MemoryDB) Remove(chatname string) error {
	if _, ok := d.chats[chatname]; !ok {
		return er.ErrNoSuchChat
	}
	delete(d.chats, chatname)
	return nil
}
