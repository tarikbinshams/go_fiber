package config

import "github.com/golang-jwt/jwt/v5"

var JwtSecret = []byte("your-secret-key")

type AuthClaims struct {
	Email  string `json:"email"`
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}
