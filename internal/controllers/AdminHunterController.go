package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"github.com/pavel-one/WebhookWatcher/internal/resources"
	"net/http"
)

func (c *AdminController) HunterCreate(w http.ResponseWriter, r *http.Request) {
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

func (c *AdminController) HunterGet(w http.ResponseWriter, r *http.Request) {
	hunter := new(models.Hunter)
	vars := mux.Vars(r)
	hunter.FindBySlug(c.DB, vars["slug"])

	if hunter.Slug == "" {
		c.Error(w, http.StatusBadRequest, errors.New(hunterErr))
		return
	}

	c.JSON(w, http.StatusOK, hunter)
}

func (c *AdminController) HunterUpdate(w http.ResponseWriter, r *http.Request) {
	hunter := new(models.Hunter)
	vars := mux.Vars(r)
	hunter.FindBySlug(c.DB, vars["slug"])

	if hunter.Id == "" {
		c.Error(w, http.StatusBadRequest, errors.New(hunterErr))
		return
	}

	var requestHunter struct {
		Ip   string `json:"ip"`
		Slug string `json:"slug"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestHunter)

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	if requestHunter.Ip == "" && requestHunter.Slug == "" {
		c.Error(w, http.StatusBadGateway, errors.New("request body is required"))
		return
	}

	if requestHunter.Ip != "" {
		hunter.Ip = requestHunter.Ip
	}

	if requestHunter.Slug != "" {
		hunter.Slug = requestHunter.Slug
	}

	if err = hunter.Update(c.DB); err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusOK, hunter)
}

func (c *AdminController) HunterDelete(w http.ResponseWriter, r *http.Request) {
	hunter := new(models.Hunter)
	vars := mux.Vars(r)
	hunter.FindBySlug(c.DB, vars["slug"])

	if hunter.Slug == "" || hunter.Id == "" {
		c.Error(w, http.StatusBadRequest, errors.New(hunterErr))
		return
	}

	if err := hunter.Delete(c.DB, hunter.Id); err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusOK, map[string]any{"message": "hunter deleted successfully", "status": http.StatusOK})
}
