package chat

import (
	"testing"

	"github.com/niakr1s/chatty-server/app/pool/events"

	"github.com/stretchr/testify/assert"
)

const username = "user"

func TestChat_AddUser(t *testing.T) {
	chat := NewChat(chatname)

	err := chat.AddUser(username)
	assert.NoError(t, err)

	err = chat.AddUser(username)
	assert.Error(t, err)

	err = chat.AddUser(username + "1")
	assert.NoError(t, err)
}

func TestChat_RemoveUser(t *testing.T) {
	chat := NewChat(chatname)

	chat.AddUser(username)

	err := chat.RemoveUser(username)
	assert.NoError(t, err)

	err = chat.RemoveUser(username)
	assert.Error(t, err)

	err = chat.RemoveUser(username + "1")
	assert.Error(t, err)
}

func TestChat_IsInChat(t *testing.T) {
	chat := NewChat(chatname)
	assert.False(t, chat.IsInChat(username))

	chat.AddUser(username)
	assert.True(t, chat.IsInChat(username))

	chat.RemoveUser(username)
	assert.False(t, chat.IsInChat(username))
}

func TestChat_notify(t *testing.T) {
	ch := make(chan events.Event)

	chat := NewChat(chatname).WithNotifyCh(ch)

	chat.AddUser(username)
	chat.AddUser(username)

	joinE := (<-ch).(*events.ChatJoinEvent)
	assert.Equal(t, joinE.Username, username)

	chat.RemoveUser(username)
	chat.RemoveUser(username)

	leaveE := (<-ch).(*events.ChatLeaveEvent)
	assert.Equal(t, leaveE.Username, username)

	select {
	case <-ch:
		assert.Fail(t, "channel should be empty")
	default:
	}
}
