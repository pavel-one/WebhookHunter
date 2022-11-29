package controllers

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"github.com/pavel-one/WebhookWatcher/internal/sqlite"
	"io"
	"log"
	"net/http"
	"strings"
)

//go:embed favicon.svg
var Favicon []byte

type RequestController struct {
	BaseController
	WriteSocketController
}

func (c *RequestController) Init(ch chan<- SocketMessage) {
	c.socketCh = ch
}

func (c *RequestController) NewRequest(w http.ResponseWriter, r *http.Request) {
	domain := helpers.GetDomainWithHost(r.Host)
	var channel string
	var hunter models.Hunter
	var RequestModel models.RequestModel
	var body []byte

	channel = r.RequestURI
	if channel == "" {
		channel = "/"
	}

	if channel == "/favicon.ico" {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(Favicon)
		return
	}

	if err := hunter.FindBySlug(domain); err != nil {
		log.Printf("[ERROR] find hunter error %s", err)
		http.Redirect(w, r, "/ui/", http.StatusMovedPermanently)
		return
	}

	if hunter.Slug == "" {
		http.Redirect(w, r, "/ui/", http.StatusMovedPermanently)
		return
	}

	err, chModel := hunter.FindChannelByPath(channel)

	if chModel.Id == 0 {
		err, chModel = hunter.CreateChannel(channel)
		c.socketCh <- EventMessage{
			Domain:  hunter.Slug,
			Channel: "root",
			Event:   "AddChannel",
			Data: ChannelResponse{
				ID:    chModel.Id,
				Path:  chModel.Path,
				Date:  humanize.Time(chModel.CreatedAt.Time),
				Count: chModel.RequestCount,
			},
		}.ToSocket()
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

	db, err := sqlite.GetDb(hunter.Slug)
	if err != nil {
		return
	}

	if err := RequestModel.Create(db); err != nil {
		c.Error(w, http.StatusBadGateway, errors.New("error save request"))
		return
	}

	counts, err := chModel.GetCounts(db)
	if err != nil {
		c.Error(w, http.StatusBadGateway, errors.New("error get counts"))
		return
	}

	c.socketCh <- EventMessage{
		Domain:  hunter.Slug,
		Channel: "root",
		Event:   "UpdateCount",
		Data:    counts,
	}.ToSocket()

	ch := chModel.GetChannel()

	c.socketCh <- EventMessage{
		Domain:  hunter.Slug,
		Channel: ch,
		Event:   "Add",
		Data:    models.FormatRequest(RequestModel),
	}.ToSocket()

	c.JSON(w, http.StatusCreated, map[string]any{
		"status": http.StatusText(http.StatusCreated),
	})
}
