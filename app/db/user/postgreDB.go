package user

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/niakr1s/chatty-server/app/models"
)

// PostgreDB ...
type PostgreDB struct {
	ctx context.Context

	pool *pgxpool.Pool
}

// CreateTableCmd ...
var CreateTableCmd = `CREATE TABLE IF NOT EXISTS users (
	"id" SERIAL PRIMARY KEY,
	"user" VARCHAR(50) NOT NULL UNIQUE,
	"email" VARCHAR(50) NOT NULL UNIQUE,
	"password_hash" VARCHAR(255) NOT NULL
);`

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
	// creating our table
	if _, err := pool.Exec(ctx, CreateTableCmd); err != nil {
		return nil, err
	}
	return &PostgreDB{ctx: ctx, pool: pool}, nil
}

// Store ...
func (d *PostgreDB) Store(u *models.FullUser) error {
	if _, err := d.pool.Exec(d.ctx, `INSERT INTO users ("user", "email", "password_hash") 
	VALUES ($1, $2, $3);`, u.UserName, u.Address, u.PasswordHash); err != nil {
		return err
	}
	return nil
}

// Update ...
func (d *PostgreDB) Update(u *models.FullUser) error {
	if _, err := d.pool.Exec(d.ctx, `UPDATE "users" 
	SET "user" = $1, "email" = $2, "password_hash" = $3 
	WHERE "user" = $1;`, u.UserName, u.Address, u.PasswordHash); err != nil {
		return err
	}
	return nil
}

// Get ...
func (d *PostgreDB) Get(username string) (models.FullUser, error) {
	res := models.FullUser{}
	row := d.pool.QueryRow(d.ctx, `SELECT "user", "email", "password_hash" FROM "users" WHERE "user" = $1;`, username)
	if err := row.Scan(&res.UserName, &res.Address, &res.PasswordHash); err != nil {
		return res, err
	}
	return res, nil
}
