package db

// ChatDB stores chats in-memory
type ChatDB interface {
	Add(chatname string) error
	Remove(chatname string) error
	GetChats() []string

	AddUser(chatname, username string) error
	RemoveUser(chatname, username string) error
	IsInChat(chatname, username string) bool
	GetUsers(chatname string) []string
}
