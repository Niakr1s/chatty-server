package store

import (
	"github.com/niakr1s/chatty-server/app/db"
)

// Store contains all databases
type Store struct {
	UserDB   db.UserDB
	ChatDB   db.ChatDB
	LoggedDB db.LoggedDB
}

// NewStore ...
func NewStore(u db.UserDB, c db.ChatDB, l db.LoggedDB) *Store {
	return &Store{UserDB: u, ChatDB: c, LoggedDB: l}
}
