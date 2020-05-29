package message

import (
	"sync"
	"time"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
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
	chat := d.chats[msg.ChatName]

	msg.ID = len(chat) + 1
	msg.Time = models.UnixTime(time.Now())

	chat = append(chat, msg)
	d.chats[msg.ChatName] = chat

	return nil
}

// GetLastNMessages ...
func (d *MemoryDB) GetLastNMessages(chatname string, n int) ([]*models.Message, error) {
	chat, ok := d.chats[chatname]

	if !ok {
		return []*models.Message{}, er.ErrNoSuchChat
	}

	if len(chat) <= n {
		return chat, nil
	}

	return chat[len(chat)-n:], nil
}
