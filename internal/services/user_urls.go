package services

import (
	"context"
	"fmt"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

type UserURLsRepository interface {
	ByUID(context.Context, int) (entities.PairArray, error)
}

type UserURLs struct {
	repo UserURLsRepository
}

func NewUserURLs(repo UserURLsRepository) *UserURLs {
	return &UserURLs{
		repo: repo,
	}
}

func (u *UserURLs) URLsForUser(ctx context.Context, uid int) (entities.PairArray, error) {
	urls, err := u.repo.ByUID(ctx, uid)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user URLs from repository: %w", err)
	}

	return urls, nil
}
