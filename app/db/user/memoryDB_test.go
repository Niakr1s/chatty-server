package userdb

import (
	"server2/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryDB_Store(t *testing.T) {
	db := NewMemoryDB()

	t.Run("simple storing", func(t *testing.T) {
		u := &models.User{Name: "user", Password: "password", PasswordHash: "passwordhash"}

		err := db.Store(u)

		assert.NoError(t, err)
	})

	t.Run("id is generated", func(t *testing.T) {
		u1 := &models.User{Name: "user1", Password: "password", PasswordHash: "passwordhash"}
		u2 := &models.User{Name: "user2", Password: "password", PasswordHash: "passwordhash"}

		db.Store(u1)
		db.Store(u2)

		assert.NotEqual(t, u1.ID, u2.ID)
	})
}

func TestMemoryDB_Get(t *testing.T) {
	db := NewMemoryDB()

	t.Run("simple get", func(t *testing.T) {
		u := &models.User{Name: "user", Password: "password", PasswordHash: "passwordhash"}

		db.Store(u)

		gotU, err := db.Get(u.ID)

		assert.NoError(t, err)
		assert.Equal(t, u, gotU)
	})
}
