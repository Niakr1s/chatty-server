package postgres

import (
	"context"
	"testing"

	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestDB_StoreAndGetAndUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parentDB, _ := NewDB(ctx, testConnStr)
	clearDB(t, parentDB)

	db := parentDB.GetUserDB()

	u := models.NewFullUser("user", "user1@example1.org", "12345")
	u.GeneratePasswordHash()
	u.ErasePassword()

	err := db.Store(u)
	assert.NoError(t, err)

	storedU, err := db.Get("user")
	assert.NoError(t, err)

	assert.Equal(t, u.UserName, storedU.UserName)
	assert.Equal(t, u.Address, storedU.Address)
	assert.Equal(t, u.PasswordHash, storedU.PasswordHash)

	u.Address = "newemail@newaddres.org"
	u.PasswordHash = "newPasswordHash"

	err = db.Update(u)
	assert.NoError(t, err)

	storedU, err = db.Get("user")
	assert.NoError(t, err)

	assert.Equal(t, u, storedU)
}
