package user

import (
	"testing"

	"github.com/niakr1s/chatty-server/app/models"

	"github.com/stretchr/testify/assert"
)

func TestMemoryDB_Store(t *testing.T) {
	db := NewMemoryDB()

	t.Run("simple storing", func(t *testing.T) {
		u := &models.User{Name: "user", Password: "password", PasswordHash: "passwordhash"}

		err := db.Store(u)
		assert.NoError(t, err)
	})

	t.Run("same user storing twice", func(t *testing.T) {
		u1 := &models.User{Name: "user1", Password: "password", PasswordHash: "passwordhash"}

		err := db.Store(u1)
		assert.NoError(t, err)

		err = db.Store(u1)
		assert.Error(t, err)
	})
}

func TestMemoryDB_Get(t *testing.T) {
	db := NewMemoryDB()

	t.Run("simple get", func(t *testing.T) {
		u := models.User{Name: "user", Password: "password", PasswordHash: "passwordhash"}

		db.Store(&u)

		gotU, err := db.Get(u.Name)

		assert.NoError(t, err)
		assert.Equal(t, u, gotU)
	})
}
