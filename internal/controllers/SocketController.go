package controllers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"strings"
)

type SocketController struct {
	BaseController
	DatabaseController
	Upgrader websocket.Upgrader
	Clients  map[string]map[*websocket.Conn]bool
}

func (c *SocketController) Init(db *sqlx.DB) {
	c.DB = db
	c.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Пропускаем любой запрос
		},
	}
}

func (c *SocketController) Test(w http.ResponseWriter, r *http.Request) {
	connection, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] Error upgrading %v", err)
		return
	}
	defer connection.Close()
	domain := strings.Split(r.Host, ".")[0]

	if c.Clients == nil {
		m := make(map[string]map[*websocket.Conn]bool)
		m[domain] = map[*websocket.Conn]bool{
			connection: true,
		}
		c.Clients = m
	} else {
		c.Clients[domain][connection] = true
	}

	defer delete(c.Clients[domain], connection)

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь с клиентом прервана
		}

		connection.WriteMessage(websocket.TextMessage, message)

		go messageHandler(message)
	}
}

func messageHandler(message []byte) {
	fmt.Println("[DEBUG]: " + string(message))
}
