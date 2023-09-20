package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ilya-rusyanov/shrinklator/internal/handlers"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
)

type PseudoAuth struct {
	log        *logger.Log
	key        string
	cookieName string
	expiration time.Duration
}

func NewPseudoAuth(log *logger.Log, key, cookieName string) *PseudoAuth {
	return &PseudoAuth{
		log:        log,
		key:        key,
		cookieName: cookieName,
		expiration: 10 * time.Minute,
	}
}

func (a *PseudoAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(a.cookieName)
		if err != nil || !a.valid(*cookie) {
			a.log.Debug("building auth cookie")
			c, err := a.buildAuthCookie()
			if err != nil {
				a.log.Error("failed to create auth cookie")
			} else {
				http.SetCookie(rw, c)
			}
		}
		next.ServeHTTP(rw, r)
	})
}

func (a *PseudoAuth) valid(cookie http.Cookie) bool {
	_, err := jwt.ParseWithClaims(cookie.Value, &handlers.Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(a.key), nil
		})

	if err != nil {
		return false
	} else {
		return true
	}
}

func (a *PseudoAuth) buildAuthCookie() (*http.Cookie, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, handlers.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.expiration)),
		},
		UserID: func() *int { val := 1; return &val }(),
	})

	tokenString, err := token.SignedString([]byte(a.key))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &http.Cookie{
		Name:  a.cookieName,
		Value: tokenString,
	}, nil
}
