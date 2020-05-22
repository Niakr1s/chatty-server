package chat

import "sync"

// DB stores chats in-memory
type DB interface {
	sync.Locker

	// if err == ErrChatAlreadyExists, returned *Chat must be valid
	Add(chatname string) (*Chat, error)

	Get(chatname string) (*Chat, error)
	Remove(chatname string) error
}
