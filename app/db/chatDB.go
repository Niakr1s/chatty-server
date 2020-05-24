package db

import (
	"sync"
)

// ChatDB stores chats in-memory
type ChatDB interface {
	sync.Locker

	// if err == ErrChatAlreadyExists, returned *Chat must be valid
	Add(chatname string) (Chat, error)

	Get(chatname string) (Chat, error)
	Remove(chatname string) error
}

// Chat ...
type Chat interface {
	ChatName() string
	AddUser(username string) error
	RemoveUser(username string) error
	IsInChat(username string) bool
}
