package models

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type Hunter struct {
	Id        string    `json:"id" db:"id"`
	Ip        string    `json:"ip" db:"ip"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (h *Hunter) Create(db *sqlx.DB) error {
	h.CreatedAt = time.Now()
	h.Id = uuid.NewString()

	_, err := db.NamedExec(`INSERT INTO hunters (id, ip, created_at) 
								VALUES (:id, :ip, :created_at)`, h)

	if err != nil {
		return err
	}

	if err := h.Find(db, h.Id); err != nil {
		return err
	}

	return nil
}

func (h *Hunter) Find(db *sqlx.DB, id string) error {
	return db.Get(h, "SELECT * FROM hunters WHERE id=$1 ORDER BY id DESC LIMIT 1", id)
}
