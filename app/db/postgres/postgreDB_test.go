package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testConnStr = "postgres://localhost:5432/users"

func TestNewPostgreDB(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := NewPostgreDB(ctx, testConnStr)

	assert.NoError(t, err)
}

func clearPostgreDB(t *testing.T, db *PostgreDB) {
	t.Helper()

	var err error
	_, err = db.pool.Exec(db.ctx, `TRUNCATE TABLE "users"`)
	assert.NoError(t, err)
	_, err = db.pool.Exec(db.ctx, `TRUNCATE TABLE "chats"`)
	assert.NoError(t, err)
	_, err = db.pool.Exec(db.ctx, `TRUNCATE TABLE "messages"`)
	assert.NoError(t, err)
}
