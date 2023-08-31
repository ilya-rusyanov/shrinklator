package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
)

var errHashing = errors.New("bad input")

type shortStorage interface {
	Put(id string, value string) error
	ByID(id string) (string, error)
}

type Shortener struct {
	storage shortStorage
}

func NewShortener(storage shortStorage) *Shortener {
	res := &Shortener{storage}
	return res
}

func (s *Shortener) Shrink(input string) (string, error) {
	hash := md5.Sum([]byte(input))
	hashStr := hex.EncodeToString(hash[:])
	err := s.storage.Put(hashStr, input)
	if err != nil {
		return "", fmt.Errorf("error storing: %w", err)
	}
	return hashStr, nil
}

func (s *Shortener) Expand(input string) (string, error) {
	url, err := s.storage.ByID(input)

	if err != nil {
		return "", fmt.Errorf("error searching: %w", err)
	}

	return url, nil
}
