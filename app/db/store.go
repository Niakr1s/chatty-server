package db

// Store contains all databases
type Store struct {
	UserDB   UserDB
	ChatDB   ChatDB
	LoggedDB LoggedDB
}

// NewStore ...
func NewStore(u UserDB, c ChatDB, l LoggedDB) *Store {
	return &Store{UserDB: u, ChatDB: c, LoggedDB: l}
}
