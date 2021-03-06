package postgres

import (
	"sync"
	"time"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/models"
)

// MessageDB is db.MessageDB impl
type MessageDB struct {
	sync.Mutex
	p *DB

	loggedDB db.LoggedDB
}

// NewMessageDB ...
func NewMessageDB(p *DB) *MessageDB {
	return &MessageDB{p: p}
}

// Post ...
// should update message ID and time
func (d *MessageDB) Post(m *models.Message) error {
	id, err := d.getUserID(m.UserName)
	if err != nil {
		id = -1
	}
	m = m.WithTime(time.Now().UTC())
	return d.postWithUserID(m, id)
}

func (d *MessageDB) postWithUserID(m *models.Message, id int) error {
	if err := d.p.pool.QueryRow(d.p.ctx, `INSERT INTO messages ("user_id", "user", "chat", "text", "time", "verified", "bot")
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		id, m.UserName, m.ChatName, m.Text, m.Time.ToSQLTimeStamp(), m.Verified, m.Bot).Scan(&m.ID); err != nil {
		return err
	}
	return nil
}

// GetLastNMessages ...
// should return empty slice even on error
func (d *MessageDB) GetLastNMessages(chatname string, n int) ([]*models.Message, error) {
	res := []*models.Message{}
	rows, err := d.p.pool.Query(d.p.ctx, `SELECT *
	FROM (SELECT id, "user", text, time, verified, bot
		FROM messages
		WHERE chat=$1
		ORDER BY id DESC
		LIMIT $2) AS q1
	ORDER BY q1.id ASC;`, chatname, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := models.NewMessage("", "", "")
		var t time.Time
		err := rows.Scan(&m.ID, &m.UserName, &m.Text, &t, &m.Verified, &m.Bot)
		if err != nil {
			return nil, err
		}
		m = m.WithTime(t)
		m.ChatName = chatname
		res = append(res, m)
	}
	return res, nil
}

func (d *MessageDB) getUserID(username string) (int, error) {
	id := 0
	if err := d.p.pool.QueryRow(d.p.ctx, `SELECT id FROM users WHERE "user"=$1`, username).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
