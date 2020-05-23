package db

import (
	"server2/app/models"
	"sync"
)

// ChatDB stores chats in-memory
type ChatDB interface {
	sync.Locker

	// if err == ErrChatAlreadyExists, returned *Chat must be valid
	Add(chatname string) (*models.Chat, error)

	Get(chatname string) (*models.Chat, error)
	Remove(chatname string) error
}
