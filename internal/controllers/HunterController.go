package controllers

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

type HunterController struct {
	DatabaseController
}

func (c *HunterController) Init(db *sqlx.DB) {
	c.DB = db
}

func (c *HunterController) Create(w http.ResponseWriter, r *http.Request) {

}
