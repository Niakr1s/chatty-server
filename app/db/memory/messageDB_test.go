package memory

import (
	"testing"
	"time"

	"github.com/niakr1s/chatty-server/app/models"

	"github.com/stretchr/testify/assert"
)

func TestMessageDB_Post(t *testing.T) {
	db := NewMessageDB()

	msg := models.NewMessage(username, text, chatname).WithTime(time.Now().UTC())
	assert.Equal(t, msg.ID, 0)

	err := db.Post(msg)

	assert.NoError(t, err)
	assert.Equal(t, msg.ID, 1)
}

func TestMessageDB_GetLastNMessages(t *testing.T) {
	db := NewMessageDB()

	msg := models.NewMessage(username, text, chatname).WithTime(time.Now().UTC())

	for i := 0; i < 10; i++ {
		db.Post(msg)
	}

	for _, i := range []int{3, 10, 16} {
		last, err := db.GetLastNMessages(chatname, i)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(last), 10)
	}
}
