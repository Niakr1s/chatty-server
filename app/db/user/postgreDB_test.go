package user

import (
	"context"
	"testing"

	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

const connStr = "user=postgres password=postgres dbname=users sslmode=disable"

// const connStr = "postgres://localhost:5432"

func clearPostgreDB(t *testing.T, db *PostgreDB) {
	t.Helper()

	_, err := db.pool.Exec(db.ctx, `TRUNCATE TABLE "users"`)
	assert.NoError(t, err)
}

func TestNewPostgreDB(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := NewPostgreDB(ctx, connStr)

	assert.NoError(t, err)
}

func TestPostgreDB_StoreAndGetAndUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, _ := NewPostgreDB(ctx, connStr)
	clearPostgreDB(t, db)

	u := models.NewFullUser("user", "user1@example1.org", "12345")
	u.GeneratePasswordHash()

	err := db.Store(&u)
	assert.NoError(t, err)

	storedU, err := db.Get("user")
	assert.NoError(t, err)

	assert.Equal(t, u.UserName, storedU.UserName)
	assert.Equal(t, u.Address, storedU.Address)
	assert.Equal(t, u.PasswordHash, storedU.PasswordHash)

	u.Address = "newemail@newaddres.org"
	u.PasswordHash = "newPasswordHash"

	err = db.Update(&u)
	assert.NoError(t, err)

	storedU, err = db.Get("user")
	assert.NoError(t, err)

	assert.Equal(t, u.UserName, storedU.UserName)
	assert.Equal(t, u.Address, storedU.Address)
	assert.Equal(t, u.PasswordHash, storedU.PasswordHash)
}
