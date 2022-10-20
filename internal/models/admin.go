package models

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type Admin struct {
	Id        uint      `db:"id"`
	Login     string    `json:"login" db:"login"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func (a *Admin) Create(db *sqlx.DB) error {
	a.CreatedAt = time.Now()
	_, err := db.NamedExec(`INSERT INTO admins (login, password, created_at) 
								VALUES (:login, :password, :created_at)`, a)

	if err != nil {
		return errors.New("failed to create admin " + err.Error())
	}

	// update model
	if err = db.Get(a, "SELECT * FROM admins ORDER BY id DESC LIMIT 1"); err != nil {
		return err
	}

	return nil
}

func (a *Admin) GetByLogin(db *sqlx.DB, login string) error {
	if err := db.Get(a, "SELECT * FROM admins WHERE login=$1 ORDER BY id DESC LIMIT 1", login); err != nil {
		return err
	}

	return nil
}

func (a *Admin) Delete(db *sqlx.DB) error {
	_, err := db.NamedExec("DELETE FROM admins WHERE id=:id", a)
	return err
}
