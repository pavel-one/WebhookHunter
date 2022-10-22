package controllers

import (
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
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
		log.Printf("[DEBUG] [%s] %s", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, ok := helpers.CheckAuthHeader(r)

		if !ok {
			helpers.WriteErrMessage(w, "auth token is missing")
			return
		}

		token, err := ParseToken(authToken)

		if err != nil {
			helpers.WriteErrMessage(w, "failed parsing token")
			return
		}

		_, ok = token.Claims.(*CustomClaims)

		if !ok || !token.Valid {
			helpers.WriteErrMessage(w, "invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
