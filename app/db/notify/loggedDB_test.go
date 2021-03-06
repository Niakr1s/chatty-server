package notify

import (
	"testing"

	"github.com/niakr1s/chatty-server/app/db/memory"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/stretchr/testify/assert"
)

func TestLoggedDB_notify(t *testing.T) {
	ch := make(chan events.Event)

	memoryDB := NewLoggedDB(memory.NewLoggedDB(), ch)

	memoryDB.Login(username)
	memoryDB.Login(username) // shouldn't fire same event twice

	loginE := (<-ch).(*events.LoginEvent)
	assert.Equal(t, loginE.UserName, username)

	memoryDB.Logout(username)
	memoryDB.Logout(username) // shouldn't fire same event twice

	logoutE := (<-ch).(*events.LogoutEvent)
	assert.Equal(t, logoutE.UserName, username)

	select {
	case <-ch:
		assert.Fail(t, "no other events expected")
	default:
	}
}
