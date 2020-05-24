package logged

import (
	"testing"
	"time"

	"github.com/niakr1s/chatty-server/app/er"

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

	memoryDB.Lock()
	memoryDB.Login(username)
	memoryDB.Unlock()

	<-time.After(time.Millisecond * 50)

	memoryDB.Lock()
	_, err := memoryDB.Get(username)
	memoryDB.Unlock()

	assert.Error(t, err)
}
