package sqlite

import (
	_ "embed"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed default.sql
var schema string

func GetDb(userSlug string) (db *sqlx.DB, err error) {
	db, err = sqlx.Open("sqlite3", "./storage/users/"+userSlug+"/"+userSlug+".sqlite3")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

func SetDefaultSchema(db *sqlx.DB) (err error) {
	_, err = db.Exec(schema)
	if err != nil {
		return err
	}

	return err
}
