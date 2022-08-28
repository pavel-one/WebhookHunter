package controllers

import (
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"github.com/pavel-one/WebhookWatcher/internal/resources"
	"net/http"
)

type HunterController struct {
	BaseController
	DatabaseController
}

func (c *HunterController) Init(db *sqlx.DB) {
	c.DB = db
}

func (c *HunterController) Create(w http.ResponseWriter, r *http.Request) {
	hunter := new(models.Hunter)
	response := new(resources.HunterResponse)

	hunter.Ip = r.RemoteAddr
	if err := hunter.Create(c.DB); err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	response.Init(hunter)
	c.JSON(w, http.StatusCreated, response)
	return
}

func (c *HunterController) Index(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": "1.0",
	})

	return
}
