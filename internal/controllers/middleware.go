package controllers

import (
	"log"
	"net/http"
)

var Handler404 = http.HandlerFunc(handler404)

func handler404(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
	log.Printf("[DEBUG] [%s] 404 %s", r.Method, r.URL.String())
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[DEBUG] [%s] %s %s", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}
