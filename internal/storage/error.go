package storage

type ErrAlreadyExists struct {
	StoredValue string
}

func (e ErrAlreadyExists) Error() string {
	return "value already exists"
}
