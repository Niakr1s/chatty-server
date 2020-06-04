package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testConnStr = "postgres://localhost:5432/users"

func newTestDB(t *testing.T) (*DB, func()) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())

	db, err := NewDB(ctx, testConnStr)
	assert.NoError(t, err)
	return db, cancel
}

func clearDB(t *testing.T, db *DB) {
	t.Helper()

	var err error
	_, err = db.pool.Exec(db.ctx, `TRUNCATE TABLE "users"`)
	assert.NoError(t, err)
	_, err = db.pool.Exec(db.ctx, `TRUNCATE TABLE "chats"`)
	assert.NoError(t, err)
	_, err = db.pool.Exec(db.ctx, `TRUNCATE TABLE "messages"`)
	assert.NoError(t, err)
}
