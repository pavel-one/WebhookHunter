package models

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Channel struct {
	Id       uint           `db:"id"`
	HunterId string         `db:"hunter_id"`
	Path     string         `db:"path"`
	Redirect sql.NullString `db:"redirect"`
}

func (c *Channel) Create(db *sqlx.DB) error {

	_, err := db.NamedExec(`INSERT INTO channels (hunter_id, path, redirect) 
								VALUES (:hunter_id, :path, :redirect)`, c)

	if err != nil {
		return errors.New("failed to create channel " + err.Error())
	}

	return nil
}

func (c *Channel) Find(db *sqlx.DB, id int64) error {
	return db.Get(c, "SELECT * FROM channels WHERE id=$1 ORDER BY id DESC LIMIT 1", id)
}

func (c *Channel) Delete(db *sqlx.DB) error {
	_, err := db.NamedExec("DELETE FROM channels WHERE id=$1", c.Id)
	return err
}

func (c *Channel) FindByPath(db *sqlx.DB, path string) error {
	return db.Get(c, "SELECT * FROM channels WHERE path=$1 LIMIT 1", path)
}
