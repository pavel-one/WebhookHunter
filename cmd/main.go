package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pavel-one/WebhookWatcher/internal/base"
	"github.com/pavel-one/WebhookWatcher/internal/controllers"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

func main() {
	fatalChan := make(chan error, 1)

	app := new(base.App)
	app.Init()

	hunterController := new(controllers.HunterController)
	hunterController.Init(app.DB)

	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": "1.0",
		})

		if err != nil {
			return
		}
	}).Methods("GET")

	app.Router.HandleFunc("/", hunterController.Create).Methods("GET")

	go app.ApiRun("8080", fatalChan)

	err := <-fatalChan
	if err != nil {
		log.Printf("[FATAL] %v", err)
		os.Exit(1)
	}
}
