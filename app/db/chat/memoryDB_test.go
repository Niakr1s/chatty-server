package chat

import (
	"server2/app/pool/events"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const chatname = "chat"

func TestMemoryDB_Add(t *testing.T) {
	chatDB := NewMemoryDB()

	c1, err := chatDB.Add("chat1")

	assert.NoError(t, err)

	c1New, err := chatDB.Add("chat1")

	assert.Error(t, err)
	assert.Equal(t, c1, c1New)
}

func TestMemoryDB_Get(t *testing.T) {
	chatDB := NewMemoryDB()

	c1, _ := chatDB.Add("chat1")
	c1New, _ := chatDB.Get("chat1")

	assert.Equal(t, c1, c1New)
}

func TestMemoryDB_Remove(t *testing.T) {
	chatDB := NewMemoryDB()

	chatDB.Add("chat1")

	err := chatDB.Remove("chat1")
	assert.NoError(t, err)

	err = chatDB.Remove("other")
	assert.Error(t, err)
}

func TestMemoryDB_notify(t *testing.T) {
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

func TestMemoryDB_ChatWithChannel(t *testing.T) {
	ch := make(chan events.Event)

	memoryDB := NewMemoryDB().WithNotifyCh(ch)

	chat, _ := memoryDB.Add(chatname)

	chat.notifyUserJoined("user", chatname, time.Now())

	joinedE := (<-ch).(*events.ChatJoinEvent)
	assert.Equal(t, joinedE.Chatname, chatname)
	assert.Equal(t, joinedE.Username, "user")
}
