package handlers

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

// Claims - JWT claims
type Claims struct {
	jwt.RegisteredClaims
	UserID *entities.UserID
}
