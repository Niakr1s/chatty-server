package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChat_AddUser(t *testing.T) {
	chat := NewChat(chatname)

	err := chat.AddUser(username)
	assert.NoError(t, err)

	err = chat.AddUser(username)
	assert.Error(t, err)

	err = chat.AddUser(username + "1")
	assert.NoError(t, err)
}

func TestChat_RemoveUser(t *testing.T) {
	chat := NewChat(chatname)

	chat.AddUser(username)

	err := chat.RemoveUser(username)
	assert.NoError(t, err)

	err = chat.RemoveUser(username)
	assert.Error(t, err)

	err = chat.RemoveUser(username + "1")
	assert.Error(t, err)
}

func TestChat_IsInChat(t *testing.T) {
	chat := NewChat(chatname)
	assert.False(t, chat.IsInChat(username))

	chat.AddUser(username)
	assert.True(t, chat.IsInChat(username))

	chat.RemoveUser(username)
	assert.False(t, chat.IsInChat(username))
}
