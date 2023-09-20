package services

import (
	"errors"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

type UserURLsRepository interface {
}

type UserURLs struct {
	repo UserURLsRepository
}

func NewUserURLs(repo UserURLsRepository) *UserURLs {
	return &UserURLs{
		repo: repo,
	}
}

func (u *UserURLs) URLsForUser(uid int) (entities.PairArray, error) {
	return entities.PairArray{}, errors.New("TODO")
}
