package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/niakr1s/chatty-server/app/internal/migrations"
	log "github.com/sirupsen/logrus"
)

// PostgreDB ...
type PostgreDB struct {
	ctx context.Context

	pool *pgxpool.Pool
}

// NewPostgreDB ...
func NewPostgreDB(ctx context.Context, connStr string) (*PostgreDB, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	log.Infof("PostgreDB: connected to %s sucsessfully", connStr)
	return &PostgreDB{ctx: ctx, pool: pool}, nil
}

// ApplyMigrations applies migrations from dir to create valid tables.
// First naiive impl, applies all migrations from folder, step by step.
func (d *PostgreDB) ApplyMigrations(migrationsDir string) error {
	migr, err := migrations.GetMigrations(migrationsDir)
	if err != nil {
		return err
	}
	for _, m := range migr {
		if _, err := d.pool.Exec(d.ctx, m.Contents); err != nil {
			return err
		}
	}
	log.Infof("PostgresDB: %d migrations from dir %s applied succesfully", len(migr), migrationsDir)
	return nil
}
