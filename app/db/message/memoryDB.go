package message

import (
	"server2/app/er"
	"server2/app/models"
	"sync"
)

// MemoryDB ...
type MemoryDB struct {
	sync.Mutex

	chats map[string][]*models.Message
}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{chats: make(map[string][]*models.Message)}
}

// Post ...
func (d *MemoryDB) Post(msg *models.Message) error {
	chat := d.chats[msg.Chat]
	msg.ID = len(chat) + 1
	chat = append(chat, msg)
	d.chats[msg.Chat] = chat
	return nil
}

// GetLastNMessages ...
func (d *MemoryDB) GetLastNMessages(chatname string, n int) ([]*models.Message, error) {
	chat, ok := d.chats[chatname]

	if !ok {
		return nil, er.ErrNoSuchChat
	}

	if len(chat) <= n {
		return chat, nil
	}

	return chat[len(chat)-n:], nil
}
