package controllers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
	"log"
	"net/http"
)

type DatabaseController struct {
	DB *sqlx.DB
}

type BaseController struct {
}

type WriteSocketController struct {
	socketCh chan<- SocketMessage
}

type UseSocketController struct {
	socketCh chan SocketMessage
}

func (c *BaseController) Error(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": http.StatusText(code),
		"error":  err.Error(),
	})

	if err != nil {
		log.Printf("[ERROR] Error response error: %v", err)
		return
	}

	log.Printf("[INFO] Error response: %v", err)
}

func (c *BaseController) JSON(w http.ResponseWriter, code int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Printf("[ERROR] Error response json encoding: %v", err)
		return
	}

	log.Printf("[INFO] JSON response: %v", response)
}

func (c *BaseController) NotFound(w http.ResponseWriter, r *http.Request) {
	helpers.Handler404(w, r)
}
