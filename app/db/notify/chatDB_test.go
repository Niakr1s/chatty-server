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

	memoryDB := NewChatDB(memory.NewChatDB(), memory.NewLoggedDB(), ch)

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

	memoryDB := NewChatDB(memory.NewChatDB(), memory.NewLoggedDB(), ch)

	err := memoryDB.Add(chatname)
	assert.NoError(t, err)
	<-ch

	memoryDB.AddUser(chatname, username)

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
	d := NewChatDB(memory.NewChatDB(), memory.NewLoggedDB(), make(chan<- events.Event))
	d.Add(chatname)
	d.AddUser(chatname, username)

	assert.NotEmpty(t, d.GetUsers(chatname))

	d.StartListeningToEvents(ch)

	ch <- &events.LogoutEvent{UserEvent: &events.UserEvent{User: models.User{UserName: username}}}

	<-time.After(time.Millisecond * 10)

	assert.Empty(t, d.GetUsers(chatname))
}

func TestChat_notify(t *testing.T) {
	ch := make(chan events.Event)

	db := NewChatDB(memory.NewChatDB(), memory.NewLoggedDB(), ch)
	db.Add(chatname)
	<-ch

	db.AddUser(chatname, username)
	db.AddUser(chatname, username)

	joinE := (<-ch).(*events.ChatJoinEvent)
	assert.Equal(t, joinE.UserName, username)

	db.RemoveUser(chatname, username)
	db.RemoveUser(chatname, username)

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

	db := NewChatDB(memory.NewChatDB(), logged, ch)
	db.Add(chatname)
	<-ch

	db.AddUser(chatname, username)

	joinE := (<-ch).(*events.ChatJoinEvent)
	assert.Equal(t, joinE.UserName, username)
	assert.Equal(t, true, joinE.Verified)
	assert.Equal(t, true, joinE.Admin)
}
