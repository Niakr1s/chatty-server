package notify

import (
	"testing"
	"time"

	"github.com/niakr1s/chatty-server/app/db/message"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestMessageDB_notify(t *testing.T) {
	ch := make(chan events.Event)

	db := NewMessageDB(message.NewMemoryDB(), ch)

	msg := models.NewMessage(username, text, chatname).WithTime(time.Now().UTC())

	for i := 0; i < 3; i++ {
		db.Post(msg)
		msgE, ok := (<-ch).(*events.MessageEvent)

		assert.True(t, ok)
		assert.Equal(t, msgE.Message, msg)
	}

	select {
	case <-ch:
		assert.Fail(t, "channel should be empty")
	default:
	}
}
