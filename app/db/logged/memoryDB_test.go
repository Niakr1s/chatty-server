package logged

import (
	"server2/app/er"
	"server2/app/pool/events"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	username = "user"
)

func TestMemoryDB_Login(t *testing.T) {
	db := NewMemoryDB()

	u1, err := db.Login("user1")

	assert.NoError(t, err)
	assert.NotEmpty(t, u1.LoginToken)

	sameU, err := db.Login("user1")

	assert.Error(t, err)

	if err == er.ErrAlreadyLogged {
		assert.Equal(t, u1, sameU)
	}

	u2, err := db.Login("user2")

	assert.NoError(t, err)
	assert.NotEmpty(t, u2.LoginToken)

	assert.NotEqual(t, u1.LoginToken, u2.LoginToken)
}

func TestMemoryDB_Get(t *testing.T) {
	db := NewMemoryDB()

	u1, _ := db.Login("user1")

	got1, err := db.Get("user1")

	assert.NoError(t, err)

	assert.Equal(t, u1, got1)
}

func TestMemoryDB_Logout(t *testing.T) {
	db := NewMemoryDB()

	db.Login("user1")

	err := db.Logout("user1")
	assert.NoError(t, err)

	err = db.Logout("user1")
	assert.Error(t, err)
}

func TestMemoryDB_StartCleanInactiveUsers(t *testing.T) {
	memoryDB := NewMemoryDB()

	memoryDB.StartCleanInactiveUsers(time.Millisecond*10, time.Millisecond*10)

	username := "user"
	memoryDB.Login(username)

	<-time.After(time.Millisecond * 50)

	memoryDB.Lock()
	_, err := memoryDB.Get(username)
	memoryDB.Unlock()

	assert.Error(t, err)
}

func TestMemoryDB_notifyLogin(t *testing.T) {
	ch := make(chan events.Event)

	memoryDB := NewMemoryDB().WithNotifyCh(ch)

	memoryDB.Login(username)

	loginE := (<-ch).(*events.LoginEvent)
	assert.Equal(t, loginE.Username, username)

	memoryDB.Logout(username)

	logoutE := (<-ch).(*events.LogoutEvent)
	assert.Equal(t, logoutE.Username, username)

	select {
	case <-ch:
		assert.Fail(t, "no other events expected")
	default:
	}
}
