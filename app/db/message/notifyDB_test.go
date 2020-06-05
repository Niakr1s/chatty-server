package message

import (
	"testing"
	"time"

	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestNotifyDB_notify(t *testing.T) {
	ch := make(chan events.Event)

	db := NewMemoryDB().WithNotifyCh(ch)

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
