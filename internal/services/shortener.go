package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
)

var errHashing = errors.New("bad input")

type shortStorage interface {
	Put(ctx context.Context, id string, value string) error
	ByID(ctx context.Context, id string) (string, error)
}

type Shortener struct {
	storage shortStorage
}

func NewShortener(storage shortStorage) *Shortener {
	res := &Shortener{storage}
	return res
}

func (s *Shortener) Shrink(ctx context.Context, input string) (string, error) {
	hash := md5.Sum([]byte(input))
	hashStr := hex.EncodeToString(hash[:])
	err := s.storage.Put(ctx, hashStr, input)
	if err != nil {
		return "", fmt.Errorf("error storing: %w", err)
	}
	return hashStr, nil
}

func (s *Shortener) Expand(ctx context.Context, input string) (string, error) {
	url, err := s.storage.ByID(ctx, input)

	if err != nil {
		return "", fmt.Errorf("error searching: %w", err)
	}

	return url, nil
}
