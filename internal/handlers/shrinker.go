package handlers

import (
	"context"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

// Shrinker - represnets shortener service
type Shrinker interface {
	Shrink(context.Context, string, *entities.UserID) (string, error)
	Expand(context.Context, string) (entities.ExpandResult, error)
}
