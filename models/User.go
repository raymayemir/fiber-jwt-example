package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID uint `json:"id"`
}

type Token struct {
	UserID string
	jwt.StandardClaims
}
