package models

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type AuthToken struct {
	Id        uint      `db:"id"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
}

func (t *AuthToken) Create(db *sqlx.DB) error {
	t.CreatedAt = time.Now()
	_, err := db.NamedExec(`INSERT INTO tokens (token, created_at) 
								VALUES (:token, :created_at)`, t)

	if err != nil {
		return errors.New("failed to create token " + err.Error())
	}

	// update model
	if err = db.Get(t, "SELECT * FROM tokens ORDER BY id DESC LIMIT 1"); err != nil {
		return err
	}

	return nil
}

func (t *AuthToken) GetByToken(db *sqlx.DB, token string) error {
	if err := db.Get(t, "SELECT * FROM tokens WHERE token=$1 ORDER BY id DESC LIMIT 1", token); err != nil {
		return err
	}

	return nil
}

func (t *AuthToken) Delete(db *sqlx.DB) error {
	_, err := db.NamedExec("DELETE FROM tokens WHERE id=:id", t)
	return err
}
