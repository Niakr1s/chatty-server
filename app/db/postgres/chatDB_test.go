package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatDB_Add(t *testing.T) {
	const chatname = "chat"
	parentDB, cancel := newTestDB(t)
	defer cancel()

	clearDB(t, parentDB)

	chatDB := NewChatDB(parentDB)

	_, err := chatDB.Add(chatname)
	assert.NoError(t, err)
	_, err = chatDB.Add(chatname)
	assert.Error(t, err)
	chat, err := chatDB.Add(chatname)
	assert.Error(t, err)

	gotChat, err := chatDB.Get(chatname)
	assert.NoError(t, err)
	assert.Equal(t, chat, gotChat)

	err = chatDB.Remove(chatname)
	assert.NoError(t, err)

	err = chatDB.Remove(chatname)
	assert.Error(t, err)
}
