package models

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
	"log"
	"os"
	"time"
)

type Hunter struct {
	Slug string `json:"slug" db:"slug"`
}

func (h *Hunter) Create() error {
	h.Slug = uuid.NewString()
	os.MkdirAll("storage/users/", os.ModePerm)

	return nil
}

func (h *Hunter) Find(db *sqlx.DB, id string) error {
	return db.Get(h, "SELECT * FROM hunters WHERE id=$1 ORDER BY id DESC LIMIT 1", id)
}

func (h *Hunter) FindBySlug(db *sqlx.DB, slug string) error {
	return db.Get(h, "SELECT * FROM hunters WHERE slug=$1 LIMIT 1", slug)
}

func (h *Hunter) CreateChannel(db *sqlx.DB, channel string) (error, Channel) {
	if h.Id == "" {
		return errors.New("hunter not exists"), Channel{}
	}

	ch := Channel{
		HunterId: h.Id,
		Path:     channel,
	}

	if err := ch.Create(db); err != nil {
		return err, Channel{}
	}

	return nil, ch
}

func (h *Hunter) FindChannelByPath(db *sqlx.DB, channel string) (error, Channel) {
	var ch Channel

	if h.Id == "" {
		return errors.New("hunter not exists"), ch
	}

	err := db.Get(&ch, "SELECT * FROM channels WHERE hunter_id=$1 and path=$2", h.Id, channel)
	if err != nil {
		return err, ch
	}

	return nil, ch
}

func (h *Hunter) Channels(db *sqlx.DB) ([]Channel, error) {
	var channels []Channel

	if h.Id == "" {
		return nil, errors.New("model not exists")
	}

	err := db.Select(&channels, `select c.id, hunter_id, c.path, c.redirect, c.created_at, count(r.id) request_count
										from channels as c
												 left join requests r on c.id = r.channel_id
										WHERE hunter_id = $1
										group by c.id`, h.Id)
	if err != nil {
		return nil, err
	}

	return channels, nil
}
