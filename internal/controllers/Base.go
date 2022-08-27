package controllers

import "github.com/jmoiron/sqlx"

type DatabaseController struct {
	DB *sqlx.DB
}
