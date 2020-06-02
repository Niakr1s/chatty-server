package user

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/niakr1s/chatty-server/app/models"
	log "github.com/sirupsen/logrus"
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
	"email_activation_token" VARCHAR(50) NOT NULL,
	"email_activated" BOOLEAN NOT NULL DEFAULT FALSE,
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
func (d *PostgreDB) Store(u models.FullUser) error {
	log.Tracef("PostgreDB: start storing %v", u)
	if _, err := d.pool.Exec(d.ctx, `INSERT INTO users 
	("user", "email", "email_activation_token", "email_activated", "password_hash") 
	VALUES ($1, $2, $3, $4, $5);`, u.UserName, u.Address, u.ActivationToken, u.Activated, u.PasswordHash); err != nil {
		return err
	}
	log.Tracef("PostgreDB: succes storing %v", u)
	return nil
}

// Update ...
func (d *PostgreDB) Update(u models.FullUser) error {
	if _, err := d.pool.Exec(d.ctx, `UPDATE "users" 
	SET "user" = $1, "email" = $2, "email_activation_token" = $3, "email_activated" = $4, "password_hash" = $5
	WHERE "user" = $1;`, u.UserName, u.Address, u.ActivationToken, u.Activated, u.PasswordHash); err != nil {
		return err
	}
	return nil
}

// Get ...
func (d *PostgreDB) Get(username string) (models.FullUser, error) {
	log.Tracef("PostgreDB: start getting %s", username)
	res := models.FullUser{}
	row := d.pool.QueryRow(d.ctx, `SELECT "user", "email", "email_activation_token", "email_activated", "password_hash" 
	FROM "users" WHERE "user" = $1;`, username)
	if err := row.Scan(&res.UserName, &res.Address, &res.ActivationToken, &res.Activated, &res.PasswordHash); err != nil {
		return res, err
	}
	log.Tracef("PostgreDB: success getting %s", username)
	return res, nil
}
