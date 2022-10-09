package controllers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type SocketMessage struct {
	Channel string
	Message string
}

type SocketController struct {
	BaseController
	DatabaseController
	Upgrader    websocket.Upgrader
	Clients     map[string][]*websocket.Conn
	MessageChan chan SocketMessage
	mu          *sync.Mutex
}

func (c *SocketController) Init(db *sqlx.DB) {
	c.DB = db
	c.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return c.checkSlug(r)
		},
	}
	c.MessageChan = make(chan SocketMessage, 10)
	c.mu = &sync.Mutex{}
}

func (c *SocketController) Test(w http.ResponseWriter, r *http.Request) {
	connection, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] Error upgrading %v", err)
		return
	}

	domain := strings.Split(r.Host, ".")[0]

	if c.Clients == nil {
		m := make(map[string][]*websocket.Conn)

		var connections []*websocket.Conn
		finConn := append(connections, connection)

		m[domain] = finConn
		c.Clients = m
	} else {
		n := append(c.Clients[domain], connection)
		c.Clients[domain] = n
	}

	if len(c.Clients[domain]) == 1 {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go c.listener(domain, wg)
		wg.Wait()
	}
}

func (c *SocketController) WorkerMessage() {
	for message := range c.MessageChan {
		log.Printf("%v connections on %v", len(c.Clients[message.Channel]), message.Channel)
		log.Println("[DEBUG] Sending message...")

		for index, connection := range c.Clients[message.Channel] {
			c.mu.Lock()
			err := connection.WriteMessage(websocket.TextMessage, []byte(message.Message))
			c.mu.Unlock()

			if err != nil {
				c.mu.Lock()
				c.Clients[message.Channel][index] = c.Clients[message.Channel][len(c.Clients[message.Channel])-1]
				c.Clients[message.Channel] = c.Clients[message.Channel][:len(c.Clients[message.Channel])-1]
				c.mu.Unlock()

				connection.Close()
				log.Printf("[ERR] Failed send message: %s", err)
				log.Printf("[WARN] connection %v closed", index+1)
				continue
			}
		}
	}
}

func (c *SocketController) listener(domain string, wg *sync.WaitGroup) {
	defer wg.Done()
	i := 0
	for {
		if len(c.Clients[domain]) != 0 {
			c.MessageChan <- SocketMessage{
				Channel: domain,
				Message: fmt.Sprintf("Send test message # %d", i),
			}

			go c.WorkerMessage()
		} else {
			log.Println("finished controller for this domain: ", domain)
			break
		}
		time.Sleep(time.Second * 2)
		i++
	}
}

func (c *SocketController) checkSlug(r *http.Request) bool {
	domain := strings.Split(r.Host, ".")

	if domain[1]+"."+domain[2] != os.Getenv("DOMAIN") {
		return false
	}

	hunter := new(models.Hunter)
	hunter.FindBySlug(c.DB, domain[0])

	if hunter.Id == "" {
		return false
	}

	return true
}
