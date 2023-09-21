package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
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
		var uid *entities.UserID

		if err != nil || !a.valid(*cookie, &uid) {
			a.log.Info("request misses auth cookie, building it")
			uid = new(entities.UserID)
			*uid = entities.UserID(rand.Intn(65535))
			c, err := a.buildAuthCookie(*uid)
			if err != nil {
				a.log.Error("failed to create auth cookie",
					zap.String("err", err.Error()))
				http.Error(rw, "cookie creation failure", http.StatusInternalServerError)
				return
			} else {
				http.SetCookie(rw, c)
			}
		} else {
			a.log.Info("request with valid auth cookie")
		}

		if uid != nil {
			a.log.Info("user id ", zap.Int("id", int(*uid)))
			ctx := context.WithValue(r.Context(), handlers.UID, *uid)
			r = r.WithContext(ctx)
		} else {
			a.log.Info("user id missing")
		}

		next.ServeHTTP(rw, r)
	})
}

func (a *PseudoAuth) valid(cookie http.Cookie, uid **entities.UserID) bool {
	claims := handlers.Claims{}
	_, err := jwt.ParseWithClaims(cookie.Value, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(a.key), nil
		})

	if err != nil {
		return false
	}

	*uid = claims.UserID

	return true
}

func (a *PseudoAuth) buildAuthCookie(uid entities.UserID) (*http.Cookie, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, handlers.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.expiration)),
		},
		UserID: &uid,
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
