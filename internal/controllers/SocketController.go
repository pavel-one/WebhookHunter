package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"log"
	"net/http"
	"os"
	"strings"
)

type SocketMessage struct {
	Domain  string
	Channel string
	Message string
}

type Connections map[*websocket.Conn]bool
type Channels map[string]Connections
type Hub map[string]Channels

type SocketController struct {
	BaseController
	DatabaseController
	Upgrader     websocket.Upgrader
	Hub          Hub
	MessageChain chan SocketMessage
}

func (c *SocketController) Init(db *sqlx.DB) {

	c.DB = db
	c.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return c.checkSlug(r)
		},
	}
	c.MessageChain = make(chan SocketMessage, 10)

	go c.WorkerMessage()
}

func (c *SocketController) Connect(w http.ResponseWriter, r *http.Request) {
	connection, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] Error upgrading %v", err)
		return
	}
	defer connection.Close()
	domain := strings.Split(r.Host, ".")[0]

	var channel string
	var channelLstn models.Channel

	if mux.Vars(r)["channel"] == "" {
		channel = "/"
		fmt.Println(channel)
	} else {
		channel = mux.Vars(r)["channel"]
	}
	var hunter models.Hunter
	hunter.FindBySlug(c.DB, domain)

	if hunter.Id == "" {
		log.Printf("[ERROR] cannot find hunter with %v slug", domain)
		return
	}

	channelLstn.Path = r.Host + channel
	channelLstn.FindByPath(c.DB, channelLstn.Path)

	if channelLstn.Id == 0 {
		channelLstn.HunterId = hunter.Id
		err = channelLstn.Create(c.DB)
		if err != nil {
			log.Printf("[ERROR] cannot create channel %v", err)
			return
		}
	}

	c.handleConnection(domain, channel, connection)
	defer c.closeConnection(connection, domain, channel)

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
	// simple add new connection
	c.Hub[domain][chName][conn] = true
}

func (c *SocketController) closeConnection(conn *websocket.Conn, domain string, chName string) {
	delete(c.Hub[domain][chName], conn)

	if len(c.Hub[domain]) == 0 {
		delete(c.Hub, domain)
		return
	}
}

func (c *SocketController) WorkerMessage() {
	for message := range c.MessageChain {
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
	if r.Host == os.Getenv("DOMAIN") {
		return false
	}

	domain := strings.Split(r.Host, ".")

	hunter := new(models.Hunter)
	hunter.FindBySlug(c.DB, domain[0])

	if hunter.Id == "" {
		return false
	}

	return true
}
