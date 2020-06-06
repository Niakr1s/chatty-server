package memory

import (
	"testing"

	"github.com/niakr1s/chatty-server/app/models"

	"github.com/stretchr/testify/assert"
)

func GenerateMockFullUser(t *testing.T) models.FullUser {
	t.Helper()
	return models.NewFullUser("user", "user@example.com", "password")
}

func TestUserDB_Store(t *testing.T) {

	t.Run("simple storing", func(t *testing.T) {
		db := NewUserDB()
		u := GenerateMockFullUser(t)

		err := db.Store(u)
		assert.NoError(t, err)
	})

	t.Run("same user storing twice", func(t *testing.T) {
		db := NewUserDB()
		u1 := GenerateMockFullUser(t)

		err := db.Store(u1)
		assert.NoError(t, err)

		err = db.Store(u1)
		assert.Error(t, err)
	})
}

func TestUserDB_Get(t *testing.T) {
	db := NewUserDB()

	t.Run("simple get", func(t *testing.T) {
		u := GenerateMockFullUser(t)

		db.Store(u)

		gotU, err := db.Get(u.UserName)

		assert.NoError(t, err)
		assert.Equal(t, u, gotU)
	})
}
