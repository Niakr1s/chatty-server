package chat

import "sync"

// MemoryDB ...
type MemoryDB struct {
	sync.Mutex

	chats map[string]*Chat
}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{chats: make(map[string]*Chat)}
}

// Add ...
func (d *MemoryDB) Add(chatname string) (*Chat, error) {
	if c, ok := d.chats[chatname]; ok {
		return c, ErrChatAlreadyExists
	}
	c := NewChat(chatname)
	d.chats[chatname] = c
	return c, nil
}

// Get ...
func (d *MemoryDB) Get(chatname string) (*Chat, error) {
	if c, ok := d.chats[chatname]; ok {
		return c, nil
	}
	return nil, ErrNoSuchChat
}

// Remove ...
func (d *MemoryDB) Remove(chatname string) error {
	if _, ok := d.chats[chatname]; !ok {
		return ErrNoSuchChat
	}
	delete(d.chats, chatname)
	return nil
}
