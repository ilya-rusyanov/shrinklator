package storage

import "context"

type Interface interface {
	Put(ctx context.Context, id, value string) error
	ByID(ctx context.Context, id string) (string, error)
}
