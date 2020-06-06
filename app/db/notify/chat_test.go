package notify

import (
	"testing"

	"github.com/niakr1s/chatty-server/app/db/memory"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/stretchr/testify/assert"
)

func TestChat_notify(t *testing.T) {
	ch := make(chan events.Event)

	chat := NewChat(memory.NewChat(chatname), memory.NewLoggedDB(), ch)

	chat.AddUser(username)
	chat.AddUser(username)

	joinE := (<-ch).(*events.ChatJoinEvent)
	assert.Equal(t, joinE.UserName, username)

	chat.RemoveUser(username)
	chat.RemoveUser(username)

	leaveE := (<-ch).(*events.ChatLeaveEvent)
	assert.Equal(t, leaveE.UserName, username)

	select {
	case <-ch:
		assert.Fail(t, "channel should be empty")
	default:
	}
}
