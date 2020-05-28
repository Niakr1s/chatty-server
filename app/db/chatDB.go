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

	GetChats() []Chat
}

// Chat ...
type Chat interface {
	sync.Locker

	ChatName() string
	AddUser(username string) error
	RemoveUser(username string) error
	IsInChat(username string) bool
}

// Chatnames ...
func Chatnames(chats []Chat) []string {
	res := make([]string, 0, len(chats))
	for _, chat := range chats {
		res = append(res, chat.ChatName())
	}
	return res
}
