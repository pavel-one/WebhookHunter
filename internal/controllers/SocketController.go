package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"github.com/pavel-one/WebhookWatcher/internal/sqlite"
	"log"
	"net/http"
)

type SocketMessage struct {
	Domain  string
	Channel string
	Message string
}

type EventMessage struct {
	Domain  string
	Channel string
	Event   string
	Data    any
}

func (e EventMessage) ToSocket() SocketMessage {
	str := make(map[string]any)
	str["message"] = ""
	str["event"] = e.Event
	str["data"] = e.Data

	marshal, err := json.Marshal(str)
	if err != nil {
		return SocketMessage{}
	}

	return SocketMessage{
		Domain:  e.Domain,
		Channel: e.Channel,
		Message: string(marshal),
	}
}

type Connections map[*websocket.Conn]bool
type Channels map[string]Connections
type Hub map[string]Channels

type SocketController struct {
	BaseController
	UseSocketController
	Upgrader websocket.Upgrader
	Hub      Hub
}

func (c *SocketController) Init(ch chan SocketMessage) {

	c.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return c.checkSlug(r)
		},
	}
	c.socketCh = ch

	go c.WorkerMessage()
}

func (c *SocketController) Connect(w http.ResponseWriter, r *http.Request) {
	var channel string
	var hunter models.Hunter

	connection, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] Error upgrading %v", err)
		return
	}
	defer connection.Close()
	domain := helpers.GetDomainWithHost(r.Host)

	channel = r.RequestURI[1:len(r.RequestURI)]
	log.Printf("[DEBUG] Socket connect [%s] channel", channel)

	if channel == "" {
		channel = "/"
	}

	err = hunter.FindBySlug(domain)

	if hunter.Slug == "" || err != nil {
		log.Printf("[WARNING] cannot find hunter with %v slug", domain)
		return
	}

	db, err := sqlite.GetDb(hunter.Slug)
	if err != nil {
		return
	}

	c.handleConnection(domain, channel, connection)

	if channel != "root" {
		var requests []models.RequestModel

		_, channelModel := hunter.FindChannelByPath("/" + channel)
		if channelModel.Id != 0 {
			requests, _ = channelModel.GetRequests(db)
		}

		c.socketCh <- EventMessage{
			domain,
			channel,
			"Load",
			models.FormatRequests(requests),
		}.ToSocket()
	}

	defer c.dropConnection(connection, domain, channel)

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil {
			return // exit
		}

		switch mt {
		case websocket.PingMessage, websocket.PongMessage:
			continue
		case websocket.TextMessage:
			messageHandler(message)
		default:
			return // exit
		}

	}
}

func (c *SocketController) handleConnection(domain string, chName string, conn *websocket.Conn) {
	// if hub not exists
	if c.Hub == nil {
		c.Hub = Hub{
			domain: Channels{
				chName: Connections{
					conn: true,
				},
			},
		}
		return
	}

	// if domain not exists
	if len(c.Hub[domain]) == 0 {
		c.Hub[domain] = Channels{
			chName: Connections{
				conn: true,
			},
		}
		return
	}

	// if channel not exists
	if _, ok := c.Hub[domain][chName]; !ok {
		c.Hub[domain][chName] = Connections{
			conn: true,
		}
		return
	}

	// simple add new connection
	c.Hub[domain][chName][conn] = true
}

func (c *SocketController) dropConnection(conn *websocket.Conn, domain string, chName string) {
	delete(c.Hub[domain][chName], conn)

	if len(c.Hub[domain]) == 0 {
		delete(c.Hub, domain)
		return
	}
}

func (c *SocketController) WorkerMessage() {
	for message := range c.socketCh {
		if len(c.Hub[message.Domain][message.Channel]) == 0 {
			continue
		}
		log.Println("[DEBUG] Sending message...")

		for connection, _ := range c.Hub[message.Domain][message.Channel] {
			err := connection.WriteMessage(websocket.TextMessage, []byte(message.Message))
			if err != nil {
				log.Printf("[ERR] Failed send message: %s", err)
				continue
			}
		}
	}
}

func messageHandler(message []byte) {
	fmt.Println("[DEBUG]: " + string(message))
}

func (c *SocketController) checkSlug(r *http.Request) bool {
	domain := helpers.GetDomainWithHost(r.Host)

	if domain == "root" {
		return true
	}

	hunter := new(models.Hunter)

	if err := hunter.FindBySlug(domain); err != nil {
		return false
	}

	if hunter.Slug == "" {
		return false
	}

	return true
}
