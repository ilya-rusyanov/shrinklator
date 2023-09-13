package handlers

import "context"

type shrinker interface {
	Shrink(context.Context, string) (string, error)
	Expand(context.Context, string) (string, error)
}
