package models

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Channel struct {
	Id           uint      `db:"id" json:"id"`
	Path         string    `db:"path" json:"path"`
	CreatedAt    null.Time `db:"created_at" json:"date"`
	RequestCount uint      `db:"request_count" json:"count"`
}

type ChannelRequestsCount struct {
	Id    uint `json:"id"`
	Count uint `json:"count"`
}

func (c *Channel) Create(db *sqlx.DB) error {
	c.CreatedAt = null.TimeFrom(time.Now())

	_, err := db.NamedExec(`INSERT INTO channels (path, created_at) 
								VALUES (:path, :created_at)`, c)

	if err != nil {
		return errors.New("failed to create channel " + err.Error())
	}

	// update model
	if err := db.Get(c, "SELECT * FROM channels ORDER BY id DESC LIMIT 1"); err != nil {
		return err
	}

	return nil
}

func (c *Channel) Find(db *sqlx.DB, id int) (err error) {
	err = db.Get(c, `SELECT channels.*, count(r.id) as request_count FROM channels 
    								LEFT JOIN requests r on channels.id = r.channel_id
    								WHERE channels.id=$1
    								GROUP BY channels.id`, id)

	return err
}

func (c *Channel) Delete(db *sqlx.DB) error {
	_, err := db.Exec("PRAGMA foreign_keys = ON; DELETE FROM channels WHERE id=$1", c.Id)
	return err
}

func (c *Channel) GetCounts(db *sqlx.DB) (model ChannelRequestsCount, err error) {
	err = db.Get(&model, `SELECT c.id, count(requests.id) as count FROM requests 
    								LEFT JOIN channels c on c.id = requests.channel_id
    								WHERE c.id=$1
    								GROUP BY c.id`, c.Id)

	return model, err
}
