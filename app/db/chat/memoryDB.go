package chat

import (
	"sync"
	"time"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/pool/events"
)

// MemoryDB ...
type MemoryDB struct {
	sync.Mutex

	chats map[string]*Chat

	notifyCh chan<- events.Event
}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{chats: make(map[string]*Chat)}
}

// WithNotifyCh ...
func (d *MemoryDB) WithNotifyCh(ch chan<- events.Event) *MemoryDB {
	d.notifyCh = ch
	return d
}

// Add ...
func (d *MemoryDB) Add(chatname string) (*Chat, error) {
	if c, ok := d.chats[chatname]; ok {
		return c, er.ErrChatAlreadyExists
	}
	c := NewChat(chatname).WithNotifyCh(d.notifyCh)
	d.chats[chatname] = c

	d.notifyChatCreated(chatname, time.Now())

	return c, nil
}

// Get ...
func (d *MemoryDB) Get(chatname string) (*Chat, error) {
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

	d.notifyChatRemoved(chatname, time.Now())

	return nil
}

func (d *MemoryDB) notifyChatCreated(chatname string, t time.Time) {
	go func() {
		if d.notifyCh != nil {
			d.notifyCh <- events.NewChatCreatedEvent(chatname, t)
		}
	}()
}

func (d *MemoryDB) notifyChatRemoved(chatname string, t time.Time) {
	go func() {
		if d.notifyCh != nil {
			d.notifyCh <- events.NewChatRemovedEvent(chatname, t)
		}
	}()
}
