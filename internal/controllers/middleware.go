package controllers

import (
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
	"log"
	"net/http"
)

var Handler404 = http.HandlerFunc(helpers.Handler404)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[DEBUG] [%s] %s", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}
