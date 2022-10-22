package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type AdminController struct {
	BaseController
	DatabaseController
}

type CustomClaims struct {
	AdminId uint   `json:"id"`
	Login   string `json:"login"`
	jwt.RegisteredClaims
}

func (c *AdminController) Init(db *sqlx.DB) {
	c.DB = db
}

func (c *AdminController) Login(w http.ResponseWriter, r *http.Request) {
	admin := new(models.Admin)
	tokenModel := new(models.AuthToken)
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

	claims := CustomClaims{
		admin.Id,
		admin.Login,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	tokenModel.Token = tokenStr

	if err = tokenModel.Create(c.DB); err != nil {
		c.Error(w, http.StatusBadRequest, err)
		return
	}

	c.JSON(w, 200, map[string]any{
		"status": "OK",
		"token":  tokenStr,
	})
}

func (c *AdminController) Test(w http.ResponseWriter, r *http.Request) {
	claims, err := c.getClaims(r)

	if err != nil {
		c.Error(w, http.StatusUnauthorized, err)
		return
	}

	c.JSON(w, 200, map[string]interface{}{
		"id":    claims.AdminId,
		"login": claims.Login,
	})
}

func ParseToken(authToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(authToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key := []byte(os.Getenv("JWT_KEY"))

		return key, nil
	})

	return token, err
}

func (c *AdminController) getClaims(r *http.Request) (*CustomClaims, error) {
	authToken, ok := helpers.CheckAuthHeader(r)

	if !ok {
		return nil, errors.New("auth token is missing")
	}

	token, err := ParseToken(authToken)

	if err != nil {
		return nil, errors.New("error parsing token")
	}

	claims, ok := token.Claims.(*CustomClaims)

	if !ok {
		return nil, errors.New("failed to get claims")
	}

	return claims, nil
}
