package chat

import (
	"testing"

	"github.com/niakr1s/chatty-server/app/events"
	"github.com/stretchr/testify/assert"
)

func TestNotifyDB_notify(t *testing.T) {
	ch := make(chan events.Event)

	memoryDB := NewMemoryDB().WithNotifyCh(ch)

	memoryDB.Add(chatname)
	memoryDB.Add(chatname) // shouldn't fire same event twice

	createdE := (<-ch).(*events.ChatCreatedEvent)
	assert.Equal(t, createdE.Chatname, chatname)

	memoryDB.Remove(chatname)
	memoryDB.Remove(chatname) // shouldn't fire same event twice

	removedE := (<-ch).(*events.ChatRemovedEvent)
	assert.Equal(t, removedE.Chatname, chatname)

	select {
	case <-ch:
		assert.Fail(t, "no other events expected")
	default:
	}
}

func TestNotifyDB_ChatNotify(t *testing.T) {
	ch := make(chan events.Event)

	memoryDB := NewMemoryDB().WithNotifyCh(ch)

	chat, err := memoryDB.Add(chatname)
	assert.NoError(t, err)
	<-ch

	chat.AddUser(username)

	e, ok := (<-ch)
	assert.True(t, ok)

	joinedE := e.(*events.ChatJoinEvent)
	assert.Equal(t, joinedE.ChatName, chatname)
	assert.Equal(t, joinedE.UserName, "user")
}
