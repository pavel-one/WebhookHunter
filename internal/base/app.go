package base

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	Root "github.com/pavel-one/WebhookWatcher"
	"log"
	"net/http"
	"time"
)

type App struct {
	Router *mux.Router
	Server *http.Server
}

func (a *App) Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("[FATAL] Not loading environment: %v", err)
	}

	if err != nil {
		log.Fatalf("[FATAL] Unable to connect to database: %v", err)
	}
	a.Router = mux.NewRouter()
}

func (a *App) ApiRun(port string, ch chan error) {
	a.Server = &http.Server{
		Handler:      a.Router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	defer a.Server.Close()

	log.Printf("[DEBUG] Running server on port %s", port)

	if err := a.Server.ListenAndServe(); err != nil {
		ch <- errors.New(fmt.Sprintf("Error server: %s", err.Error()))
	}
}

func (a *App) Close() {
	if err := a.Server.Close(); err != nil {
		log.Fatalf("[FATAL] Unable to close server: %v", err)
		return
	}
}

func (a *App) GET(path string, method func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, method).Methods("GET")
}

func (a *App) POST(path string, method func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, method).Methods("POST")
}

func (a *App) Prefix(prefix string, f http.HandlerFunc) {
	a.Router.PathPrefix(prefix).HandlerFunc(f)
}

func (a *App) Static(prefix string, root string) {
	handler := Root.AssetHandler(prefix, root)
	a.Router.PathPrefix(prefix).Handler(handler)
}
