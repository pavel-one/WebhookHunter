package models

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
	"log"
	"time"
)

type Hunter struct {
	Id        string    `json:"id" db:"id"`
	Ip        string    `json:"ip" db:"ip"`
	Slug      string    `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (h *Hunter) Create(db *sqlx.DB) error {
	h.CreatedAt = time.Now()
	h.Id = uuid.NewString()

	s := ""
	i := 3
	// Just for fun
	for {
		count := 0
		s = helpers.RandString(i)
		err := db.Get(&count, "SELECT count(*) FROM hunters WHERE slug=$1", s)
		if err != nil {
			return err
		}

		if count == 0 {
			break
		}

		i++

		if i > 100 {
			break
		}
	}

	h.Slug = s

	tx := db.MustBegin()
	tx.MustExec(`
WITH hunter as (
    INSERT INTO hunters (id, ip, created_at, slug)
        VALUES ($1, $2, $3, $4)
        RETURNING id
)
INSERT INTO channels (hunter_id, path)
VALUES (
        (select hunter.id from hunter),
        '/') 
`, h.Id, h.Ip, h.CreatedAt, h.Slug)
	err := tx.Commit()
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return errors.New("failed create hunter")
	}

	if err := h.Find(db, h.Id); err != nil {
		return err
	}

	return nil
}

func (h *Hunter) Find(db *sqlx.DB, id string) error {
	return db.Get(h, "SELECT * FROM hunters WHERE id=$1 ORDER BY id DESC LIMIT 1", id)
}

func (h *Hunter) All(db *sqlx.DB) ([]Hunter, error) {
	var hunters []Hunter
	// TODO: maybe in future add cache here?
	if err := db.Select(&hunters, "SELECT * FROM hunters"); err != nil {
		return nil, err
	}

	return hunters, nil
}

func (h *Hunter) FindBySlug(db *sqlx.DB, slug string) error {
	return db.Get(h, "SELECT * FROM hunters WHERE slug=$1 LIMIT 1", slug)
}

func (h *Hunter) CreateChannel(db *sqlx.DB, channel string) (Channel, error) {
	if h.Id == "" {
		return Channel{}, errors.New("hunter not exists")
	}

	ch := Channel{
		HunterId: h.Id,
		Path:     channel,
	}

	if err := ch.Create(db); err != nil {
		return Channel{}, err
	}

	return ch, nil
}

func (h *Hunter) FindChannelByPath(db *sqlx.DB, channel string) (Channel, error) {
	var ch Channel

	if h.Id == "" {
		return ch, errors.New("hunter not exists")
	}

	err := db.Get(&ch, "SELECT * FROM channels WHERE hunter_id=$1 and path=$2", h.Id, channel)
	if err != nil {
		return ch, err
	}

	return ch, nil
}

func (h *Hunter) Update(db *sqlx.DB) error {
	_, err := db.NamedExec("UPDATE hunters SET ip=:ip, slug=:slug WHERE id=:id", h)

	if err != nil {
		return err
	}

	if err = h.Find(db, h.Id); err != nil {
		return err
	}

	return nil
}

func (h *Hunter) Delete(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM hunters WHERE id=$1", id)
	return err
}

func (h *Hunter) Channels(db *sqlx.DB) ([]Channel, error) {
	var channels []Channel

	if h.Id == "" {
		return nil, errors.New("model not exists")
	}

	err := db.Select(&channels, `SELECT * FROM channels WHERE hunter_id=$1`, h.Id)
	if err != nil {
		return nil, err
	}

	return channels, nil
}
