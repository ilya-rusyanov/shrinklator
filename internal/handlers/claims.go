package handlers

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID *entities.UserID
}
