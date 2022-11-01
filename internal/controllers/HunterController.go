package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"github.com/pavel-one/WebhookWatcher/internal/resources"
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

func (c *HunterController) GetChannels(w http.ResponseWriter, r *http.Request) {
	var hunter models.Hunter
	domainArr := strings.Split(r.Host, ".")
	if len(domainArr) < 3 {
		c.NotFound(w, r)
		return
	}

	domain := domainArr[0]

	if err := hunter.FindBySlug(c.DB, domain); err != nil {
		c.Error(w, http.StatusNotFound, errors.New("not found"))
		return
	}

	if hunter.Id == "" {
		c.NotFound(w, r)
		return
	}

	channels, err := hunter.Channels(c.DB)
	if err != nil {
		c.NotFound(w, r)
		return
	}

	c.JSON(w, http.StatusOK, channels)
}

func (c *HunterController) Index(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": "1.0",
	})

}
