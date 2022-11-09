package sqlite

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
)

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
	file, err := os.Open("./internal/sqlite/default.sql")
	if err != nil {
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(b))
	if err != nil {
		return err
	}

	return err
}
