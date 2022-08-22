package app

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type App struct {
	DB     *sqlx.DB
	Router *mux.Router
}

func (a *App) Init() error {
	return nil
}

func (a *App) Run() error {
	return nil
}
