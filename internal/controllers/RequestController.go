package controllers

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"log"
	"net/http"
	"strings"
)

type RequestController struct {
	BaseController
	DatabaseController
}

func (c *RequestController) Init(db *sqlx.DB) {
	c.DB = db
}

func (c *RequestController) NewRequest(w http.ResponseWriter, r *http.Request) {
	domain := strings.Split(r.Host, ".")[0]
	var channel string
	var hunter models.Hunter

	channel = mux.Vars(r)["channel"]
	if channel == "" {
		channel = "/"
	}

	if err := hunter.FindBySlug(c.DB, domain); err != nil {
		log.Printf("[ERROR] find hunter error %s", err)
		c.Error(w, http.StatusBadRequest, errors.New("hunter not found"))
		return
	}

	if hunter.Id == "" {
		c.Error(w, http.StatusBadRequest, errors.New("hunter not found"))
		return
	}

	err, chModel := hunter.FindChannelByPath(c.DB, channel)

	if chModel.Id == 0 {
		err, chModel = hunter.CreateChannel(c.DB, channel)
	}

	if err != nil {
		c.Error(w, http.StatusBadGateway, errors.New("channel not found"))
		return
	}

	// TODO: Create request

	c.JSON(w, http.StatusCreated, map[string]any{
		"status": "created",
	})
}
