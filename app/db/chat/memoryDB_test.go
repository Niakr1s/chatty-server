package chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryDB_Add(t *testing.T) {
	chatDB := NewMemoryDB()

	c1, err := chatDB.Add("chat1")

	assert.NoError(t, err)

	c1New, err := chatDB.Add("chat1")

	assert.Error(t, err)
	assert.Equal(t, c1, c1New)
}

func TestMemoryDB_Get(t *testing.T) {
	chatDB := NewMemoryDB()

	c1, _ := chatDB.Add("chat1")
	c1New, _ := chatDB.Get("chat1")

	assert.Equal(t, c1, c1New)
}

func TestMemoryDB_Remove(t *testing.T) {
	chatDB := NewMemoryDB()

	chatDB.Add("chat1")

	err := chatDB.Remove("chat1")
	assert.NoError(t, err)

	err = chatDB.Remove("other")
	assert.Error(t, err)
}
