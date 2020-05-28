package db

// Store contains all databases
type Store struct {
	UserDB    UserDB
	ChatDB    ChatDB
	LoggedDB  LoggedDB
	MessageDB MessageDB
}

// NewStore ...
func NewStore(u UserDB, c ChatDB, l LoggedDB, m MessageDB) *Store {
	return &Store{UserDB: u, ChatDB: c, LoggedDB: l, MessageDB: m}
}
