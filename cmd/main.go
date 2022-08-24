package main

import (
	"fmt"
	"github.com/pavel-one/WebhookWatcher/internal/base"
	"net/http"
)

func main() {
	app := new(base.App)
	app.Init()

	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	}).Methods("GET")

	app.Run("8080")
}
