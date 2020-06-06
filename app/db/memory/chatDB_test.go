package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatDB_Add(t *testing.T) {
	chatDB := NewChatDB()

	err := chatDB.Add("chat1")
	assert.NoError(t, err)

	err = chatDB.Add("chat1")
	assert.Error(t, err)
}

func TestChatDB_Remove(t *testing.T) {
	chatDB := NewChatDB()

	chatDB.Add("chat1")

	err := chatDB.Remove("chat1")
	assert.NoError(t, err)

	err = chatDB.Remove("other")
	assert.Error(t, err)
}

func TestChat_AddUser(t *testing.T) {
	chatDB := NewChatDB()
	chatDB.Add(chatname)

	err := chatDB.AddUser(chatname, username)
	assert.NoError(t, err)

	err = chatDB.AddUser(chatname, username)
	assert.Error(t, err)

	err = chatDB.AddUser(chatname, username+"1")
	assert.NoError(t, err)
}

func TestChat_RemoveUser(t *testing.T) {
	chatDB := NewChatDB()
	chatDB.Add(chatname)

	chatDB.AddUser(chatname, username)

	err := chatDB.RemoveUser(chatname, username)
	assert.NoError(t, err)

	err = chatDB.RemoveUser(chatname, username)
	assert.Error(t, err)

	err = chatDB.RemoveUser(chatname, username+"1")
	assert.Error(t, err)
}

func TestChat_IsInChat(t *testing.T) {
	chatDB := NewChatDB()
	chatDB.Add(chatname)

	assert.False(t, chatDB.IsInChat(chatname, username))

	chatDB.AddUser(chatname, username)
	assert.True(t, chatDB.IsInChat(chatname, username))

	chatDB.RemoveUser(chatname, username)
	assert.False(t, chatDB.IsInChat(chatname, username))
}
