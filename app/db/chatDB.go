package db

import (
	"server2/app/db/chat"
	"sync"
)

// ChatDB stores chats in-memory
type ChatDB interface {
	sync.Locker

	// if err == ErrChatAlreadyExists, returned *Chat must be valid
	Add(chatname string) (*chat.Chat, error)

	Get(chatname string) (*chat.Chat, error)
	Remove(chatname string) error
}
