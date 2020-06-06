package postgres

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// I made everything in one test, because of concurrency.
// TODO: search internet and refactor it via mock or something else.
func TestChatDB_Add(t *testing.T) {
	const chatname = "chat"
	parentDB, cancel := newTestDB(t)
	defer cancel()

	clearDB(t, parentDB)

	chatDB := NewChatDB(parentDB)

	err := chatDB.Add(chatname)
	assert.NoError(t, err)
	err = chatDB.Add(chatname)
	assert.Error(t, err)
	err = chatDB.Add(chatname)
	assert.Error(t, err)

	err = chatDB.Remove(chatname)
	assert.NoError(t, err)

	err = chatDB.Remove(chatname)
	assert.Error(t, err)

	clearDB(t, parentDB)

	for i := 0; i < 10; i++ {
		cn := fmt.Sprintf("chat%d", i)
		chatDB.Add(cn)
		// removing chats from memory
		chatDB.ChatDB.Remove(cn)
	}
	assert.Empty(t, chatDB.ChatDB.GetChats())

	chatDB.LoadChatsFromPostgres()

	assert.Len(t, chatDB.GetChats(), 10)
}
