package adminApi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"net/http"
)

type AdminChannelController struct {
	SubController
}

func (c *AdminChannelController) Init(db *sqlx.DB) {
	c.DB = db
}

func (c *AdminChannelController) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var reqChannel struct {
		Channel string `json:"channel"`
	}
	hunter, err := c.checkHunter(vars["slug"])

	if err != nil {
		c.Error(w, http.StatusNotFound, err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&reqChannel)

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	if reqChannel.Channel == "" {
		c.Error(w, http.StatusBadRequest, errors.New("path for channel is required"))
		return
	}

	channel, _ := hunter.FindChannelByPath(c.DB, "/"+reqChannel.Channel)

	if channel.Id != 0 {
		c.Error(w, http.StatusBadRequest, errors.New("this channel already exists"))
		return
	}

	channel, err = hunter.CreateChannel(c.DB, "/"+reqChannel.Channel)

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusCreated, channel)
}

func (c *AdminChannelController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hunter, err := c.checkHunter(vars["slug"])

	if err != nil {
		c.Error(w, http.StatusNotFound, err)
	}

	channel, code, err := c.checkChannel(hunter, vars["channel"])

	if err != nil {
		c.Error(w, code, err)
		return
	}

	c.JSON(w, code, channel)
}

func (c *AdminChannelController) GetAllByHunter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hunter, err := c.checkHunter(vars["slug"])

	if err != nil {
		c.Error(w, http.StatusNotFound, err)
	}

	channels, err := hunter.AllChannelsByHunter(c.DB)

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusOK, channels)
}

func (c *AdminChannelController) GetAll(w http.ResponseWriter, r *http.Request) {
	channels, err := new(models.Channel).All(c.DB)
	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusOK, channels)
}

func (c *AdminChannelController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var requestChannel struct {
		Path     string `json:"path"`
		Redirect string `json:"redirect"`
	}
	hunter, err := c.checkHunter(vars["slug"])

	if err != nil {
		c.Error(w, http.StatusNotFound, err)
		return
	}

	channel, code, err := c.checkChannel(hunter, vars["channel"])

	if err != nil {
		c.Error(w, code, err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&requestChannel)

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	if requestChannel.Path == "" && requestChannel.Redirect == "" {
		c.Error(w, http.StatusBadRequest, errors.New("field path is required"))
		return
	}

	channel.Path = "/" + requestChannel.Path
	channel.Redirect = &sql.NullString{String: requestChannel.Redirect, Valid: true}
	if err = channel.Update(c.DB); err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, code, channel)
}

func (c *AdminChannelController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hunter, err := c.checkHunter(vars["slug"])

	if err != nil {
		c.Error(w, http.StatusNotFound, err)
		return
	}

	channel, code, err := c.checkChannel(hunter, vars["channel"])

	if err != nil {
		c.Error(w, code, err)
		return
	}

	if err = channel.Delete(c.DB, channel.Id); err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, code, map[string]any{"status": "OK", "message": "hunter deleted successfully"})
}
