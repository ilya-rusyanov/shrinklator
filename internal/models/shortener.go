package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

type Shortener struct {
	storage *storage.Storage
}

func (s *Shortener) Shrink(input string) string {
	hash := md5.Sum([]byte(input))
	hashStr := hex.EncodeToString(hash[:])
	s.storage.Put(hashStr, input)
	return hashStr
}

var errHashing = errors.New("bad input")

func (s *Shortener) Expand(input string) (string, error) {
	url, err := s.storage.ByID(input)

	if err != nil {
		return "", fmt.Errorf("error searching: %w", err)
	}

	return url, nil
}

func New(storage *storage.Storage) *Shortener {
	res := &Shortener{storage}
	return res
}
