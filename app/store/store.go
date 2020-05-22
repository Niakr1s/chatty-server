package store

import userdb "server2/app/db/user"

// Store contains all databases
type Store struct {
	userDB userdb.DB
}

// NewStore ...
func NewStore(u userdb.DB) *Store {
	return &Store{userDB: u}
}

// NewMemoryStore ...
func NewMemoryStore() *Store {
	userDB := userdb.NewMemoryDB()
	return &Store{userDB: userDB}
}
