package AdminControllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"io"
	"net/http"
	"strconv"
)

type AdminRequestController struct {
	SubController
}

func (c *AdminRequestController) Init(db *sqlx.DB) {
	c.DB = db
}

func (c *AdminRequestController) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	request := new(models.RequestModel)
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.Error(w, http.StatusBadGateway, errors.New("error read body"))
		return
	}

	if !json.Valid(body) {
		c.Error(w, http.StatusBadGateway, errors.New("json invalid"))
		return
	}

	headers, err := json.Marshal(r.Header)
	if err != nil {
		c.Error(w, http.StatusBadGateway, errors.New("error read headers"))
		return
	}

	q, err := json.Marshal(r.URL.Query())
	if err != nil {
		c.Error(w, http.StatusBadGateway, errors.New("failed parse query"))
		return
	}

	request.ChannelID = channel.Id
	request.Request = body
	request.Headers = headers
	request.Path = r.URL.Path
	request.Query = q

	if err = request.Create(c.DB); err != nil {
		c.Error(w, http.StatusBadGateway, errors.New("error save request"))
		return
	}

	c.JSON(w, http.StatusCreated, request)
}

func (c *AdminRequestController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	request := new(models.RequestModel)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	request.Find(c.DB, uint(id))

	if request.ID == 0 {
		c.Error(w, http.StatusNotFound, errors.New("request not found"))
		return
	}

	c.JSON(w, http.StatusOK, request)
}

func (c *AdminRequestController) GetAll(w http.ResponseWriter, r *http.Request) {
	requests, err := new(models.RequestModel).All(c.DB)
	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusOK, requests)
}

func (c *AdminRequestController) GetAllByChannel(w http.ResponseWriter, r *http.Request) {
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

	requests, err := new(models.RequestModel).AllByChannelId(c.DB, channel.Id)

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusOK, requests)
}

func (c AdminRequestController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var requestReq struct {
		Request json.RawMessage `json:"request"`
		Headers json.RawMessage `json:"headers"`
		Query   json.RawMessage `json:"query"`
		Path    string          `json:"path"`
	}
	hunter, err := c.checkHunter(vars["slug"])

	if err != nil {
		c.Error(w, http.StatusNotFound, err)
		return
	}

	_, code, err := c.checkChannel(hunter, vars["channel"])

	if err != nil {
		c.Error(w, code, err)
		return
	}

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	request := new(models.RequestModel)
	request.Find(c.DB, uint(id))

	if request.ID == 0 {
		c.Error(w, http.StatusNotFound, errors.New("request not found"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&requestReq)

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	if requestReq.Request == nil && requestReq.Headers == nil && requestReq.Query == nil && requestReq.Path == "" {
		c.Error(w, http.StatusBadRequest, errors.New("request body is required"))
		return
	}

	if requestReq.Request != nil {
		request.Request = requestReq.Request
	}

	if requestReq.Headers != nil {
		request.Headers = requestReq.Headers
	}

	if requestReq.Query != nil {
		request.Query = requestReq.Query
	}

	if requestReq.Path != "" {
		request.Path = requestReq.Path
	}

	if err = request.Update(c.DB); err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusOK, request)
}

func (c *AdminRequestController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	request := new(models.RequestModel)
	request.Find(c.DB, uint(id))

	if request.ID == 0 {
		c.Error(w, http.StatusNotFound, errors.New("request not found"))
		return
	}

	if err = request.Delete(c.DB, uint(id)); err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, http.StatusOK, map[string]any{"status": "OK", "message": "request deleted successfully"})
}
