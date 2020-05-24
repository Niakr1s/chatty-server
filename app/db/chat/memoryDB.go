package chat

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/pool/events"
)

// MemoryDB ...
type MemoryDB struct {
	sync.Mutex

	chats map[string]db.Chat
}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{chats: make(map[string]db.Chat)}
}

// WithNotifyCh ...
func (d *MemoryDB) WithNotifyCh(ch chan<- events.Event) *NotifyDB {
	return NewNotifyDB(d, ch)
}

// Add ...
func (d *MemoryDB) Add(chatname string) (db.Chat, error) {
	if c, ok := d.chats[chatname]; ok {
		return c, er.ErrChatAlreadyExists
	}
	c := NewMemoryChat(chatname)
	d.chats[chatname] = c

	return c, nil
}

// Get ...
func (d *MemoryDB) Get(chatname string) (db.Chat, error) {
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
