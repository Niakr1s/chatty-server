package store

import (
	"server2/app/db/chat"
	"server2/app/db/logged"
	"server2/app/db/user"
)

// Store contains all databases
type Store struct {
	UserDB   user.DB
	ChatDB   chat.DB
	LoggedDB logged.DB
}

// NewStore ...
func NewStore(u user.DB, c chat.DB, l logged.DB) *Store {
	return &Store{UserDB: u, ChatDB: c, LoggedDB: l}
}

// NewMemoryStore ...
func NewMemoryStore() *Store {
	userDB := user.NewMemoryDB()
	chatDB := chat.NewMemoryDB()
	loggedDB := logged.NewMemoryDB()
	return NewStore(userDB, chatDB, loggedDB)
}
