package models

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
	"time"
)

type RequestModel struct {
	ID        uint            `json:"id" db:"id"`
	ChannelID uint            `json:"channel_id" db:"channel_id"`
	Request   json.RawMessage `json:"request" db:"request"`
	Headers   json.RawMessage `json:"headers" db:"headers"`
	Query     json.RawMessage `json:"query" db:"query"`
	Path      string          `json:"path" db:"path"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

func (m *RequestModel) Create(db *sqlx.DB) error {
	m.CreatedAt = time.Now()

	if len(m.Request) > 0 {
		m.Request = helpers.TrimJson(m.Request)
	}

	if len(m.Headers) > 0 {
		m.Headers = helpers.TrimJson(m.Headers)
	}

	_, err := db.NamedExec(`INSERT INTO requests (request, created_at, channel_id, headers, path, query) 
						VALUES (:request, :created_at, :channel_id, :headers, :path, :query)`, m)
	if err != nil {
		return err
	}

	// update model
	if err := db.Get(m, "SELECT * FROM requests ORDER BY id DESC LIMIT 1"); err != nil {
		return err
	}

	return nil
}

func (m *RequestModel) Find(db *sqlx.DB, id uint) error {
	return db.Get(m, "SELECT * FROM requests WHERE id=$1 ORDER BY id DESC LIMIT 1", id)
}

func (m *RequestModel) All(db *sqlx.DB) ([]RequestModel, error) {
	var requests []RequestModel
	if err := db.Select(&requests, "SELECT * FROM requests"); err != nil {
		return nil, err
	}

	return requests, nil
}

func (m *RequestModel) AllByChannelId(db *sqlx.DB, channelId uint) ([]RequestModel, error) {
	var requests []RequestModel
	if err := db.Select(&requests, "SELECT * FROM requests WHERE channel_id=$1", channelId); err != nil {
		return nil, err
	}

	return requests, nil
}

func (m *RequestModel) Update(db *sqlx.DB) error {
	_, err := db.NamedExec("UPDATE requests SET request=:request, headers=:headers, path=:path, query=:query WHERE id=:id", m)

	if err != nil {
		return err
	}

	if err = m.Find(db, m.ID); err != nil {
		return err
	}

	return nil
}

func (m *RequestModel) Delete(db *sqlx.DB, id uint) error {
	_, err := db.Exec("DELETE FROM requests WHERE id=$1", id)
	return err
}
