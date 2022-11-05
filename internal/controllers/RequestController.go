package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"io"
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
	var RequestModel models.RequestModel
	var body []byte

	channel = strings.ReplaceAll(r.RequestURI, "/request/", "/")
	if channel == "" {
		channel = "/"
	}

	if err := hunter.FindBySlug(c.DB, domain); err != nil {
		log.Printf("[ERROR] find hunter error %s", err)
		c.Error(w, http.StatusBadRequest, errors.New("hunter not found"))
		return
	}

	if hunter.Id == "" {
		c.Error(w, http.StatusNotFound, errors.New("hunter not found"))
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

	headers, err := json.Marshal(r.Header)
	if err != nil {
		c.Error(w, http.StatusBadGateway, errors.New("error read headers"))
		return
	}

	typeRequest := strings.Split(r.Header.Get("Content-Type"), ";")[0]

	switch typeRequest {
	case "application/json":
		body, err = io.ReadAll(r.Body)
		if err != nil {
			c.Error(w, http.StatusBadGateway, errors.New("error read body"))
			return
		}

		if !json.Valid(body) {
			c.Error(w, http.StatusBadGateway, errors.New("json invalid"))
			return
		}

		break
	case "text/plain", "text/html", "application/xml":
		body, err = io.ReadAll(r.Body)
		if err != nil {
			c.Error(w, http.StatusBadGateway, errors.New("error read body"))
			return
		}

		body, err = json.Marshal(map[string]string{
			"type": "simple",
			"text": string(body),
		})
		break
	case "multipart/form-data":
		if err := r.ParseMultipartForm(512 << 20); err != nil {
			c.Error(w, http.StatusBadGateway, errors.New("error form data"))
			return
		}
		body, err = json.Marshal(r.MultipartForm)
		if err != nil {
			c.Error(w, http.StatusBadGateway, errors.New("error form data"))
			return
		}
		break
	case "": //not content type
		body = []byte{123, 125}
		break
	default: //other content type
		c.Error(w, http.StatusBadGateway, errors.New(fmt.Sprintf("not support content-type: %s", typeRequest)))
		return
	}

	q, err := json.Marshal(r.URL.Query())
	if err != nil {
		c.Error(w, http.StatusBadGateway, errors.New("failed parse query"))
		return
	}

	RequestModel.ChannelID = chModel.Id
	RequestModel.Request = body
	RequestModel.Headers = headers
	RequestModel.Path = r.URL.Path
	RequestModel.Query = q

	if err := RequestModel.Create(c.DB); err != nil {
		c.Error(w, http.StatusBadGateway, errors.New("error save request"))
		return
	}

	c.JSON(w, http.StatusCreated, map[string]any{
		"status": http.StatusText(http.StatusCreated),
	})
}
