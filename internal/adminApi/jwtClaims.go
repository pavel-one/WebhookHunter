package adminApi

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	AdminId uint   `json:"id"`
	Login   string `json:"login"`
	jwt.RegisteredClaims
}
