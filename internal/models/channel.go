package models

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Channel struct {
	Id           uint            `db:"id" json:"id"`
	HunterId     string          `db:"hunter_id" json:"hunter_id"`
	Path         string          `db:"path" json:"path"`
	Redirect     *sql.NullString `db:"redirect" json:"-"`
	CreatedAt    null.Time       `db:"created_at" json:"date"`
	RequestCount uint            `db:"request_count" json:"count"`
}

func (c *Channel) Create(db *sqlx.DB) error {
	c.CreatedAt = null.TimeFrom(time.Now())

	_, err := db.NamedExec(`INSERT INTO channels (hunter_id, path, redirect, created_at) 
								VALUES (:hunter_id, :path, :redirect, :created_at)`, c)

	if err != nil {
		return errors.New("failed to create channel " + err.Error())
	}

	// update model
	if err := db.Get(c, "SELECT * FROM channels ORDER BY id DESC LIMIT 1"); err != nil {
		return err
	}

	return nil
}

func (c *Channel) Find(db *sqlx.DB, id int64) error {
	var err error
	err = db.Get(c, "SELECT * FROM channels WHERE id=$1 ORDER BY id DESC LIMIT 1", id)
	err = db.Get(c.RequestCount, "SELECT count() FROM requests WHERE channel_id=$1", id)

	return err
}

func (c *Channel) Delete(db *sqlx.DB) error {
	_, err := db.NamedExec("DELETE FROM channels WHERE id=$1", c.Id)
	return err
}
