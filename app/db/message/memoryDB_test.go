package message

import (
	"server2/app/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	username = "user"
	text     = "text"
	chatname = "chat"
)

func TestMemoryDB_Post(t *testing.T) {
	db := NewMemoryDB()

	msg := models.NewMessage(username, text, chatname).WithTime(time.Now())
	assert.Equal(t, msg.ID, 0)

	err := db.Post(msg)

	assert.NoError(t, err)
	assert.Equal(t, msg.ID, 1)
}

func TestMemoryDB_GetLastNMessages(t *testing.T) {
	db := NewMemoryDB()

	msg := models.NewMessage(username, text, chatname).WithTime(time.Now())

	for i := 0; i < 10; i++ {
		db.Post(msg)
	}

	for _, i := range []int{3, 10, 16} {
		last, err := db.GetLastNMessages(chatname, i)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(last), 10)
	}
}
