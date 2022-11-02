package middlewars

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/adminApi/Auth"
	"github.com/pavel-one/WebhookWatcher/internal/models"
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

func CheckToken(next http.Handler, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		admin := new(models.Admin)
		authToken, ok := Auth.CheckAuthHeader(r)

		if !ok {
			Auth.WriteErrMessage(w, "auth token is missing")
			return
		}

		tokenModel := new(models.AuthToken)
		tokenModel.GetByToken(db, authToken)

		if tokenModel.Id == 0 {
			Auth.WriteErrMessage(w, "token not exists")
			return
		}

		token, err := Auth.ParseToken(authToken)

		if err != nil {
			if err.(*jwt.ValidationError).Errors == 16 {
				tokenModel.Delete(db)
			}

			Auth.WriteErrMessage(w, err.Error())
			return
		}

		claims, ok := token.Claims.(*Auth.CustomClaims)

		if !ok || !token.Valid {
			Auth.WriteErrMessage(w, "invalid token")
			return
		}

		if err = admin.GetByLogin(db, claims.Login); err != nil {
			Auth.WriteErrMessage(w, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
