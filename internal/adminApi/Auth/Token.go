package Auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"strings"
)

type CustomClaims struct {
	AdminId uint   `json:"id"`
	Login   string `json:"login"`
	jwt.RegisteredClaims
}

//check auth header if it exists

func CheckAuthHeader(r *http.Request) (string, bool) {
	authHeader := strings.Split(r.Header.Get("Authorization"), " ")

	if len(authHeader) < 2 || authHeader[1] == "" {
		return "", false
	}

	return authHeader[1], true
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

// write error message in response writer for middleware

func WriteErrMessage(w http.ResponseWriter, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": errorMsg,
	})
}

func GetClaims(r *http.Request) (*CustomClaims, error) {
	authToken, ok := CheckAuthHeader(r)

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
