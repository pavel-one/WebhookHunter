package models

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type Hunter struct {
	Id        string    `json:"id" db:"id"`
	Ip        string    `json:"ip" db:"ip"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (h *Hunter) Create(db *sqlx.DB) (result sql.Result, err error) {
	h.Id = uuid.NewString()

	return db.NamedExec(`INSERT INTO hunters (id, ip, created_at) 
								VALUES (:id, :ip, :created_at)`, h)
}
