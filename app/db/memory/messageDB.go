package memory

import (
	"sync"
	"time"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
)

// MessageDB ...
type MessageDB struct {
	sync.Mutex

	chats map[string][]*models.Message

	notifyCh chan<- events.Event
}

// NewMessageDB ...
func NewMessageDB() *MessageDB {
	return &MessageDB{chats: make(map[string][]*models.Message)}
}

// Post ...
func (d *MessageDB) Post(msg *models.Message) error {
	d.Lock()
	defer d.Unlock()

	return d.post(msg)
}

// GetLastNMessages ...
func (d *MessageDB) GetLastNMessages(chatname string, n int) ([]*models.Message, error) {
	d.Lock()
	defer d.Unlock()

	return d.getLastNMessages(chatname, n)
}

// concurrency-unsafe

// Post ...
func (d *MessageDB) post(msg *models.Message) error {
	chat := d.chats[msg.ChatName]

	msg.ID = len(chat) + 1
	msg.Time = models.UnixTime(time.Now().UTC())

	chat = append(chat, msg)
	d.chats[msg.ChatName] = chat

	return nil
}

// GetLastNMessages ...
func (d *MessageDB) getLastNMessages(chatname string, n int) ([]*models.Message, error) {
	chat, ok := d.chats[chatname]

	if !ok {
		return []*models.Message{}, er.ErrNoSuchChat
	}

	if len(chat) <= n {
		return chat, nil
	}

	return chat[len(chat)-n:], nil
}
