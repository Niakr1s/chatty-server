package notify

import (
	"testing"

	"github.com/niakr1s/chatty-server/app/db/memory"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
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

func TestChat_NotifyJoinChatWithUserStatus(t *testing.T) {
	logged := memory.NewLoggedDB()
	lu, _ := logged.Login(username)
	lu.UserStatus = models.UserStatus{Admin: true, Verified: true}
	logged.Update(lu)

	ch := make(chan events.Event)
	chat := NewChat(memory.NewChat(chatname), logged, ch)
	chat.AddUser(username)

	joinE := (<-ch).(*events.ChatJoinEvent)
	assert.Equal(t, joinE.UserName, username)
	assert.Equal(t, true, joinE.Verified)
	assert.Equal(t, true, joinE.Admin)
}
