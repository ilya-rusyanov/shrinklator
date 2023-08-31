package storage

type Interface interface {
	Put(id, value string) error
	ByID(id string) (string, error)
}
