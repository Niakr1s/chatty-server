package postgres

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/db/chat"
)

// ChatDB wraps chat.MemoryDB to persist our chats
type ChatDB struct {
	sync.Mutex
	*chat.MemoryDB
	p *DB
}

// NewChatDB constructs ChatDB with parent DB
func NewChatDB(p *DB) *ChatDB {
	return &ChatDB{p: p, MemoryDB: chat.NewMemoryDB()}
}

// LoadChatsFromPostgres loads chats from postgres
func (d *ChatDB) LoadChatsFromPostgres() {
	rows, err := d.p.pool.Query(d.p.ctx, `SELECT "chat" FROM "chats" ORDER BY "chat" ASC;`)
	if err != nil {
		return
	}
	defer rows.Close()

	var s string
	for rows.Next() {
		err := rows.Scan(&s)
		if err != nil || s == "" {
			return
		}
		d.Add(s)
	}
}

// Add ...
// if err == ErrChatAlreadyExists, returned *Chat must be valid
func (d *ChatDB) Add(chatname string) (db.Chat, error) {
	res, err := d.MemoryDB.Add(chatname)
	// we are trusting MemoryDB, so if added - adding it into postgres
	if err == nil {
		if _, err := d.p.pool.Exec(d.p.ctx, `INSERT INTO "chats" ("chat") VALUES ($1) ON CONFLICT DO NOTHING;`, chatname); err != nil {
			return res, err
		}
	}
	return res, err
}

// Get ...
func (d *ChatDB) Get(chatname string) (db.Chat, error) {
	return d.MemoryDB.Get(chatname)
}

// Remove ...
func (d *ChatDB) Remove(chatname string) error {
	err := d.MemoryDB.Remove(chatname)
	if err == nil {
		if _, err := d.p.pool.Exec(d.p.ctx, `DELETE FROM "chats" WHERE "chat"=$1;`, chatname); err != nil {
			return err
		}
	}
	return err
}

// GetChats ...
func (d *ChatDB) GetChats() []db.Chat {
	return d.MemoryDB.GetChats()
}
