package logged

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryDB_Login(t *testing.T) {
	db := NewMemoryDB()

	u1, err := db.Login("user1")

	assert.NoError(t, err)
	assert.NotEmpty(t, u1.loginToken)

	sameU, err := db.Login("user1")

	assert.Error(t, err)

	if err == ErrAlreadyLogged {
		assert.Equal(t, u1, sameU)
	}

	u2, err := db.Login("user2")

	assert.NoError(t, err)
	assert.NotEmpty(t, u2.loginToken)

	assert.NotEqual(t, u1.loginToken, u2.loginToken)
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
