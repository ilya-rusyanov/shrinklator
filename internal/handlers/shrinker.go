package handlers

import (
	"context"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

type shrinker interface {
	Shrink(context.Context, string, *entities.UserID) (string, error)
	Expand(context.Context, string) (string, error)
}
