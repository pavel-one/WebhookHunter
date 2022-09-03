package models

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
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

	_, err := db.NamedExec(`INSERT INTO hunters (id, ip, slug, created_at) 
								VALUES (:id, :ip, :slug, :created_at)`, h)

	if err != nil {
		return errors.New("failed to create hunter")
	}

	if err := h.Find(db, h.Id); err != nil {
		return err
	}

	return nil
}

func (h *Hunter) Find(db *sqlx.DB, id string) error {
	return db.Get(h, "SELECT * FROM hunters WHERE id=$1 ORDER BY id DESC LIMIT 1", id)
}

func (h *Hunter) FindBySlug(db *sqlx.DB, slug string) error {
	return db.Get(h, "SELECT * FROM hunters WHERE slug=$1 LIMIT 1", slug)
}