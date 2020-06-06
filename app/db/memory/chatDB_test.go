package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatDB_Add(t *testing.T) {
	chatDB := NewChatDB()

	c1, err := chatDB.Add("chat1")

	assert.NoError(t, err)

	c1New, err := chatDB.Add("chat1")

	assert.Error(t, err)
	assert.Equal(t, c1, c1New)
}

func TestChatDB_Get(t *testing.T) {
	chatDB := NewChatDB()

	c1, _ := chatDB.Add("chat1")
	c1New, _ := chatDB.Get("chat1")

	assert.Equal(t, c1, c1New)
}

func TestChatDB_Remove(t *testing.T) {
	chatDB := NewChatDB()

	chatDB.Add("chat1")

	err := chatDB.Remove("chat1")
	assert.NoError(t, err)

	err = chatDB.Remove("other")
	assert.Error(t, err)
}
