package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

// UserURLsRepository stores URLs by user
type UserURLsRepository interface {
	ByUID(context.Context, entities.UserID) (entities.PairArray, error)
	Delete(context.Context, entities.DeleteRequest) error
}

// UserURLs - usecase for URLs of a user
type UserURLs struct {
	repo    UserURLsRepository
	delErrs chan<- error
	delChan chan entities.DeleteRequest
}

// NewUserURLs constructs UserURLs object
func NewUserURLs(ctx context.Context, repo UserURLsRepository, bufsize int) (
	service *UserURLs,
	deleteErrors <-chan error,
) {
	de := make(chan error)
	deleteErrors = de
	service = &UserURLs{
		repo:    repo,
		delErrs: de,
		delChan: make(chan entities.DeleteRequest, bufsize),
	}
	go service.flushDelete(ctx)
	return
}

// URLsForUser returns URLs of a single user
func (u *UserURLs) URLsForUser(ctx context.Context, uid entities.UserID) (entities.PairArray, error) {
	urls, err := u.repo.ByUID(ctx, uid)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user URLs from repository: %w", err)
	}

	return urls, nil
}

// Delete equeues deletion of a user URL
func (u *UserURLs) Delete(ctx context.Context, req entities.DeleteRequest) error {
	u.delChan <- req
	return nil
}

// Close finalizes object
func (u *UserURLs) Close() {
	close(u.delChan)
}

func (u *UserURLs) flushDelete(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	var requests entities.DeleteRequest

loop:
	for {
		select {
		case req, ok := <-u.delChan:
			if ok {
				requests = append(requests, req...)
			}
		case <-ticker.C:
			err := u.feed(ctx, requests)
			if err != nil {
				u.delErrs <- err
			}
			requests = nil
		case <-ctx.Done():
			tempCtx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
			defer cancel()
			err := u.feed(tempCtx, requests)
			if err != nil {
				u.delErrs <- err
			}
			requests = nil
			break loop
		}
	}
	close(u.delErrs)
}

func (u *UserURLs) feed(ctx context.Context, requests entities.DeleteRequest) error {
	if len(requests) == 0 {
		return nil
	}
	err := u.repo.Delete(ctx, requests)
	return err
}
