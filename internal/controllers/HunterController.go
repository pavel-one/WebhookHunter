package controllers

import (
	"encoding/json"
	"errors"
	"github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
	"github.com/pavel-one/WebhookWatcher/internal/sqlite"
	"net/http"
	"strconv"
	"strings"

	"github.com/pavel-one/WebhookWatcher/internal/models"
	"github.com/pavel-one/WebhookWatcher/internal/resources"
)

type HunterController struct {
	BaseController
	WriteSocketController
}

type ChannelResponse struct {
	ID    uint   `json:"id"`
	Path  string `json:"path"`
	Date  string `json:"date"`
	Count uint   `json:"count"`
}

func (c *HunterController) Init(ch chan<- SocketMessage) {
	c.socketCh = ch
}

func (c *HunterController) Create(w http.ResponseWriter, r *http.Request) {
	hunter := new(models.Hunter)
	response := new(resources.HunterResponse)

	if err := hunter.Create(); err != nil {
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

	if err := hunter.FindBySlug(request.Uri); err != nil {
		if request.Uri == "localhost" {
			hunter.Slug = "localhost"
			if err := hunter.CreateWithName(); err != nil {
				c.Error(w, http.StatusBadGateway, errors.New("failed create"))
				return
			}
		} else {
			c.Error(w, http.StatusNotFound, errors.New("not found"))
			return
		}
	}

	if hunter.Slug == "" {
		c.Error(w, http.StatusNotFound, errors.New("not found"))
		return
	}

	c.JSON(w, http.StatusOK, map[string]any{
		"status": http.StatusText(http.StatusOK),
	})

}

func (c *HunterController) GetChannels(w http.ResponseWriter, r *http.Request) {
	var hunter models.Hunter
	var response []ChannelResponse

	domain := helpers.GetDomainWithHost(r.Host)

	if err := hunter.FindBySlug(domain); err != nil {
		c.Error(w, http.StatusNotFound, errors.New("not found"))
		return
	}

	if hunter.Slug == "" {
		c.NotFound(w, r)
		return
	}

	channels, err := hunter.Channels()
	if err != nil {
		c.NotFound(w, r)
		return
	}

	for _, channel := range channels {
		response = append(response, ChannelResponse{
			ID:    channel.Id,
			Path:  channel.Path,
			Date:  humanize.Time(channel.CreatedAt.Time),
			Count: channel.RequestCount,
		})
	}

	c.JSON(w, http.StatusOK, response)
}

func (c *HunterController) DropChannel(w http.ResponseWriter, r *http.Request) {
	var hunter models.Hunter
	var channel models.Channel

	domainArr := strings.Split(r.Host, ".")
	if len(domainArr) < 3 {
		c.NotFound(w, r)
		return
	}
	domain := strings.Split(r.Host, ".")[0]

	if err := hunter.FindBySlug(domain); err != nil {
		c.NotFound(w, r)
		return
	}

	db, err := sqlite.GetDb(hunter.Slug)
	if err != nil {
		c.NotFound(w, r)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err := channel.Find(db, id); err != nil {
		c.NotFound(w, r)
		return
	}

	if channel.Id == 0 {
		c.NotFound(w, r)
		return
	}

	if err := channel.Delete(db); err != nil {
		c.NotFound(w, r)
		return
	}

	c.socketCh <- EventMessage{
		Domain:  hunter.Slug,
		Channel: "root",
		Event:   "DropChannel",
		Data: ChannelResponse{
			ID:    channel.Id,
			Path:  channel.Path,
			Date:  humanize.Time(channel.CreatedAt.Time),
			Count: channel.RequestCount,
		},
	}.ToSocket()

	c.JSON(w, http.StatusOK, map[string]any{
		"status": http.StatusText(http.StatusOK),
	})
	return
}

func (c *HunterController) Index(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": "1.0",
	})

}
