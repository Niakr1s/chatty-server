package postgres

import (
	"github.com/niakr1s/chatty-server/app/models"
	log "github.com/sirupsen/logrus"
)

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
