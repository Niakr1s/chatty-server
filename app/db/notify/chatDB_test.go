package notify

import (
	"testing"
	"time"

	"github.com/niakr1s/chatty-server/app/db/memory"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestNotifyDB_notify(t *testing.T) {
	ch := make(chan events.Event)

	memoryDB := NewChatDB(memory.NewChatDB(), ch)

	memoryDB.Add(chatname)
	memoryDB.Add(chatname) // shouldn't fire same event twice

	createdE := (<-ch).(*events.ChatCreatedEvent)
	assert.Equal(t, createdE.ChatName, chatname)

	memoryDB.Remove(chatname)
	memoryDB.Remove(chatname) // shouldn't fire same event twice

	removedE := (<-ch).(*events.ChatRemovedEvent)
	assert.Equal(t, removedE.ChatName, chatname)

	select {
	case <-ch:
		assert.Fail(t, "no other events expected")
	default:
	}
}

func TestChatDB_Notify(t *testing.T) {
	ch := make(chan events.Event)

	memoryDB := NewChatDB(memory.NewChatDB(), ch)

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

func TestChatDB_StartListeningToEvents(t *testing.T) {
	const (
		chatname = "chat"
		username = "user"
	)

	ch := make(chan events.Event)
	d := NewChatDB(memory.NewChatDB(), make(chan<- events.Event))
	chat, _ := d.Add(chatname)
	chat.AddUser(username)

	assert.NotEmpty(t, chat.GetUsers())

	d.StartListeningToEvents(ch)

	ch <- &events.LogoutEvent{UserEvent: &events.UserEvent{User: models.User{UserName: username}}}

	<-time.After(time.Millisecond * 10)

	assert.Empty(t, chat.GetUsers())
}
