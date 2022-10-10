package controllers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"strings"
)

type SocketMessage struct {
	Channel string
	Message string
}

type Connections map[*websocket.Conn]bool
type Hub map[string]Connections

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
			// TODO: Проверять наличие такого хантера по slug и url
			return true // Пропускаем любой запрос
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

	c.handleConnection(domain, connection)
	defer c.closeConnection(connection, domain)

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

func (c *SocketController) handleConnection(domain string, conn *websocket.Conn) {
	// if hub not exists
	if c.Hub == nil {
		c.Hub = Hub{
			domain: Connections{
				conn: true,
			},
		}
		return
	}

	// if domain not exists
	if len(c.Hub[domain]) == 0 {
		c.Hub[domain] = Connections{
			conn: true,
		}
		return
	}

	// simple add new connection
	c.Hub[domain][conn] = true
}

func (c *SocketController) closeConnection(conn *websocket.Conn, domain string) {
	delete(c.Hub[domain], conn)

	if len(c.Hub[domain]) == 0 {
		delete(c.Hub, domain)
		return
	}
}

func (c *SocketController) WorkerMessage() {
	for message := range c.MessageChain {
		if len(c.Hub[message.Channel]) == 0 {
			continue
		}
		log.Println("[DEBUG] Sending message...")

		for connection, _ := range c.Hub[message.Channel] {
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
