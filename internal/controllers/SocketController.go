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

func (c *SocketController) Test(w http.ResponseWriter, r *http.Request) {
	connection, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] Error upgrading %v", err)
		return
	}
	defer connection.Close()
	domain := strings.Split(r.Host, ".")[0]

	if c.Hub == nil {
		c.Hub = Hub{
			domain: Connections{
				connection: true,
			},
		}
	} else {
		c.Hub[domain][connection] = true
	}

	defer c.closeConnection(connection, domain)
	//
	//if c.Clients == nil {
	//	m := make(map[string]map[*websocket.Conn]bool)
	//	m[domain] = map[*websocket.Conn]bool{
	//		connection: true,
	//	}
	//	c.Clients = m
	//} else {
	//	c.Clients[domain][connection] = true
	//}
	//
	//defer delete(c.Clients[domain], connection)

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			return // Выходим из функции, если клиент пытается закрыть соединение или связь с клиентом прервана
		}

		messageHandler(message)
	}
}

func (c *SocketController) closeConnection(conn *websocket.Conn, domain string) {
	if len(c.Hub[domain]) == 1 {
		c.Hub[domain] = nil
		return
	}

	delete(c.Hub[domain], conn)
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
