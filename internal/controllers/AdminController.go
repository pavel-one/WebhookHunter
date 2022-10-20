package controllers

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

type AdminController struct {
	BaseController
	DatabaseController
}

func (c *AdminController) Init(db *sqlx.DB) {
	c.DB = db
}

func (c *AdminController) Login(w http.ResponseWriter, r *http.Request) {
	admin := new(models.Admin)
	err := json.NewDecoder(r.Body).Decode(admin)
	p := admin.Password

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	admin.GetByLogin(c.DB, admin.Login)

	if admin.Id == 0 {
		err = errors.New("admin not found")
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(p))

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    admin.Id,
		"login": admin.Login,
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	c.JSON(w, 200, map[string]any{
		"status": "OK",
		"token":  tokenStr,
	})
}
