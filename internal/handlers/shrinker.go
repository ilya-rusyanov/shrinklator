package handlers

type shrinker interface {
	Shrink(string) (string, error)
	Expand(string) (string, error)
}
