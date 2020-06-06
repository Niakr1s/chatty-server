package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

// DB ...
type DB struct {
	ctx context.Context

	pool *pgxpool.Pool
}

// NewDB ...
func NewDB(ctx context.Context, connStr string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	log.Infof("PostgresDB: connected to %s sucsessfully", connStr)

	res := &DB{ctx: ctx, pool: pool}
	res.ApplyMigrations()

	return res, nil
}

// ApplyMigrations applies migrations from dir to create valid tables.
// First naiive impl, applies all migrations from folder, step by step.
func (d *DB) ApplyMigrations() {
	for _, m := range migrations {
		if i, err := d.pool.Exec(d.ctx, m); err != nil {
			log.Infof("PostgresDB: couldn't apply %d migration, check it", i)
		}
	}
	log.Infof("PostgresDB: %d migrations applied succesfully", len(migrations))
}
