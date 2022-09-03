package controllers

import (
	"encoding/json"
	"errors"
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

func (c *HunterController) Check(w http.ResponseWriter, r *http.Request) {
	var hunter models.Hunter
	var request struct {
		Uri string `json:"uri"`
	}

	json.NewDecoder(r.Body).Decode(&request)

	if request.Uri == "" {
		c.Error(w, http.StatusBadRequest, errors.New("uri is required"))
		return
	}

	if err := hunter.FindBySlug(c.DB, request.Uri); err != nil {
		c.Error(w, http.StatusNotFound, errors.New("not found"))
		return
	}

	if hunter.Id == "" {
		c.Error(w, http.StatusNotFound, errors.New("not found"))
		return
	}

	c.JSON(w, http.StatusOK, map[string]any{
		"status": "OK",
	})

}

func (c *HunterController) Index(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": "1.0",
	})

	return
}
