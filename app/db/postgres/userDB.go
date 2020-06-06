package postgres

import (
	"github.com/niakr1s/chatty-server/app/models"
	log "github.com/sirupsen/logrus"
)

// UserDB is a sub db, that uses pool from DB.
type UserDB struct {
	p *DB
}

// NewUserDB constructs ChatDB with parent DB
func NewUserDB(p *DB) *UserDB {
	return &UserDB{p: p}
}

// Store ...
func (d *UserDB) Store(u models.FullUser) error {
	log.Tracef("PostgresDB: start storing %v", u)
	if _, err := d.p.pool.Exec(d.p.ctx, `INSERT INTO users 
	("user", "email", "email_activation_token", "verified", "password_hash", "admin") 
	VALUES ($1, $2, $3, $4, $5, $6);`,
		u.UserName, u.Address, u.ActivationToken, u.Verified, u.PasswordHash, u.Admin); err != nil {
		return err
	}
	log.Tracef("PostgresDB: succes storing %v", u)
	return nil
}

// Update ...
func (d *UserDB) Update(u models.FullUser) error {
	if _, err := d.p.pool.Exec(d.p.ctx, `UPDATE "users" 
	SET "user" = $1, "email" = $2, "email_activation_token" = $3, "verified" = $4, "password_hash" = $5, "admin" = $6
	WHERE "user" = $1;`, u.UserName, u.Address, u.ActivationToken, u.Verified, u.PasswordHash, u.Admin); err != nil {
		return err
	}
	return nil
}

// Get ...
func (d *UserDB) Get(username string) (models.FullUser, error) {
	log.Tracef("PostgresDB: start getting %s", username)
	res := models.FullUser{}
	row := d.p.pool.QueryRow(d.p.ctx, `SELECT "user", "email", "email_activation_token", "verified", "password_hash" , "admin"
	FROM "users" WHERE "user" = $1;`, username)
	if err := row.Scan(&res.UserName, &res.Address, &res.ActivationToken, &res.Verified, &res.PasswordHash, &res.Admin); err != nil {
		return res, err
	}
	log.Tracef("PostgresDB: success getting %s", username)
	return res, nil
}
