package message

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/niakr1s/chatty-server/app/pool/events"
)

// MemoryDB ...
type MemoryDB struct {
	sync.Mutex

	chats map[string][]*models.Message

	notifyCh chan<- events.Event
}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{chats: make(map[string][]*models.Message)}
}

// WithNotifyCh ...
func (d *MemoryDB) WithNotifyCh(ch chan<- events.Event) *NotifyDB {
	return NewNotifyDB(d, ch)
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
