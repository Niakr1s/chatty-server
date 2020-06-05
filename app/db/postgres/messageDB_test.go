package postgres

import (
	"testing"
	"time"

	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestMessageDB(t *testing.T) {
	parentDB, cancel := newTestDB(t)
	defer cancel()

	messageDB := NewMessageDB(parentDB)

	// for lastnmessages
	messageDB.postWithUserID(models.NewMessage("user", "hello", "main"), 1)
	messageDB.postWithUserID(models.NewMessage("user", "hello", "main"), 1)
	messageDB.postWithUserID(models.NewMessage("user", "hello", "main"), 1)

	const sz = 4

	messages := make([]*models.Message, sz)
	for i := range messages {
		messages[i] = models.NewMessage("user", "hello", "main")
		m := messages[i]
		err := messageDB.postWithUserID(messages[i], 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, m.ID)
		assert.NotEmpty(t, m.Time)
	}

	last, err := messageDB.GetLastNMessages("main", sz)
	assert.NoError(t, err)
	assert.NotEmpty(t, last)
	if !assert.Len(t, last, sz) {
		t.Fatal()
	}

	for i := 0; i < sz; i++ {
		m := messages[i]
		lastM := last[i]
		assert.Equal(t, m.ID, lastM.ID)
		assert.Equal(t, m.UserName, lastM.UserName)
		assert.Equal(t, m.Text, lastM.Text)
		assert.Equal(t, m.ChatName, lastM.ChatName)

		assert.Equal(t, time.Time(m.Time).Year(), time.Time(lastM.Time).Year())
		assert.Equal(t, time.Time(m.Time).Month(), time.Time(lastM.Time).Month())
		assert.Equal(t, time.Time(m.Time).Day(), time.Time(lastM.Time).Day())
		assert.Equal(t, time.Time(m.Time).Hour(), time.Time(lastM.Time).Hour())
		assert.Equal(t, time.Time(m.Time).Minute(), time.Time(lastM.Time).Minute())
		assert.Equal(t, time.Time(m.Time).Second(), time.Time(lastM.Time).Second())
	}

	last, err = messageDB.GetLastNMessages("some other chat", sz)
	assert.NoError(t, err)
	assert.Empty(t, last)
}
