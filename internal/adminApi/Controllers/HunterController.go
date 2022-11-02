package AdminControllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"github.com/pavel-one/WebhookWatcher/internal/resources"
	"net/http"
)

type AdminHunterController struct {
	SubController
}

func (c *AdminHunterController) Init(db *sqlx.DB) {
	c.DB = db
}

func (c *AdminHunterController) Create(w http.ResponseWriter, r *http.Request) {
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

func (c *AdminHunterController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hunter, err := c.checkHunter(vars["slug"])

	if err != nil {
		c.Error(w, http.StatusNotFound, err)
		return
	}

	c.JSON(w, http.StatusOK, hunter)
}

func (c *AdminHunterController) GetAll(w http.ResponseWriter, r *http.Request) {
	hunters, err := new(models.Hunter).All(c.DB)
	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusOK, hunters)
}

func (c *AdminHunterController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hunter, err := c.checkHunter(vars["slug"])

	if err != nil {
		c.Error(w, http.StatusNotFound, err)
		return
	}

	var requestHunter struct {
		Ip   string `json:"ip"`
		Slug string `json:"slug"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestHunter)

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	if requestHunter.Ip == "" && requestHunter.Slug == "" {
		c.Error(w, http.StatusBadRequest, errors.New("request body is required"))
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

func (c *AdminHunterController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hunter, err := c.checkHunter(vars["slug"])

	if err != nil {
		c.Error(w, http.StatusNotFound, err)
		return
	}

	if err = hunter.Delete(c.DB, hunter.Id); err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusOK, map[string]any{"message": "hunter deleted successfully", "status": "OK"})
}
