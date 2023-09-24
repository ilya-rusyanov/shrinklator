package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

type UserURLsRepository interface {
	ByUID(context.Context, entities.UserID) (entities.PairArray, error)
}

type UserURLs struct {
	repo UserURLsRepository
}

func NewUserURLs(repo UserURLsRepository) *UserURLs {
	return &UserURLs{
		repo: repo,
	}
}

func (u *UserURLs) URLsForUser(ctx context.Context, uid entities.UserID) (entities.PairArray, error) {
	urls, err := u.repo.ByUID(ctx, uid)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user URLs from repository: %w", err)
	}

	return urls, nil
}

func (u *UserURLs) Delete(ctx context.Context, req entities.DeleteRequest) error {
	return errors.New("TODO")
}
