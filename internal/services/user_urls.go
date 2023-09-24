package services

import (
	"context"
	"fmt"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

type UserURLsRepository interface {
	ByUID(context.Context, entities.UserID) (entities.PairArray, error)
	Delete(context.Context, entities.DeleteRequest) error
}

type UserURLs struct {
	repo    UserURLsRepository
	delErrs chan<- error
}

func NewUserURLs(repo UserURLsRepository) (service *UserURLs,
	deleteErrors <-chan error) {
	de := make(chan error)
	deleteErrors = de
	service = &UserURLs{
		repo:    repo,
		delErrs: de,
	}
	return
}

func (u *UserURLs) URLsForUser(ctx context.Context, uid entities.UserID) (entities.PairArray, error) {
	urls, err := u.repo.ByUID(ctx, uid)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user URLs from repository: %w", err)
	}

	return urls, nil
}

func (u *UserURLs) Delete(ctx context.Context, req entities.DeleteRequest) error {
	go func() {
		err := u.repo.Delete(ctx, req)

		if err != nil {
			u.delErrs <- err
		}
	}()

	return nil
}
