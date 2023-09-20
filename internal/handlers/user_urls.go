package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
)

var noUserID = errors.New("cookie does not contain user ID")

type URLsService interface {
	URLsForUser(int) (entities.PairArray, error)
}

type UserURLs struct {
	log        *logger.Log
	service    URLsService
	key        string
	cookieName string
	baseURL    string
}

func NewUserURLs(log *logger.Log, service URLsService, tokenKey,
	cookieName, baseURL string) *UserURLs {
	return &UserURLs{
		log:        log,
		service:    service,
		key:        tokenKey,
		cookieName: cookieName,
		baseURL:    baseURL,
	}
}

func (u *UserURLs) Handler(rw http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(u.cookieName)
	if err != nil {
		http.Error(rw, fmt.Sprintf("cannot get cookie from request: %q", err.Error()),
			http.StatusInternalServerError)
		return
	}

	id, err := u.getUID(cookie)

	if err != nil {
		if errors.Is(err, noUserID) {
			http.Error(rw, "user ID is expected", http.StatusUnauthorized)
			return
		} else {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}

	urls, err := u.service.URLsForUser(id)

	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to fetch URLs: %q", err.Error()),
			http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	for i := range urls {
		urls[i].Short = u.baseURL + "/" + urls[i].Short
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(urls); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *UserURLs) getUID(cookie *http.Cookie) (id int, err error) {
	claims := &Claims{}
	_, err = jwt.ParseWithClaims(cookie.Value, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(u.key), nil
		})

	if err != nil {
		return -1, fmt.Errorf("cannot parse token: %w", err)
	}

	if claims.UserID == nil {
		return -1, noUserID
	}

	return *claims.UserID, nil
}
