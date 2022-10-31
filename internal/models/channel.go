package models

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Channel struct {
	Id       uint            `json:"id" db:"id"`
	HunterId string          `json:"hunter_id" db:"hunter_id"`
	Path     string          `json:"path" db:"path"`
	Redirect *sql.NullString `json:"redirect" db:"redirect"`
}

func (c *Channel) Create(db *sqlx.DB) error {

	_, err := db.NamedExec(`INSERT INTO channels (hunter_id, path, redirect) 
								VALUES (:hunter_id, :path, :redirect)`, c)

	if err != nil {
		return errors.New("failed to create channel " + err.Error())
	}

	// update model
	if err = db.Get(c, "SELECT * FROM channels ORDER BY id DESC LIMIT 1"); err != nil {
		return err
	}

	return nil
}

func (c *Channel) Find(db *sqlx.DB, id uint) error {
	return db.Get(c, "SELECT * FROM channels WHERE id=$1 ORDER BY id DESC LIMIT 1", id)
}

func (c *Channel) All(db *sqlx.DB) ([]Channel, error) {
	var channels []Channel
	if err := db.Select(&channels, "SELECT * FROM channels"); err != nil {
		return nil, err
	}

	return channels, nil
}

func (c *Channel) Update(db *sqlx.DB) error {
	_, err := db.NamedExec("UPDATE channels SET path=:path, redirect=:redirect WHERE id=:id", c)

	if err != nil {
		return err
	}

	if err = c.Find(db, c.Id); err != nil {
		return err
	}

	return nil
}

func (c *Channel) Delete(db *sqlx.DB, id uint) error {
	_, err := db.Exec("DELETE FROM channels WHERE id=$1", id)
	return err
}
