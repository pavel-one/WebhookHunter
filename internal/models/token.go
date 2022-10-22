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
