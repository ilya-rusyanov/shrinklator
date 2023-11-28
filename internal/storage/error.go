package storage

// ErrAlreadyExists represents not found error
type ErrAlreadyExists struct {
	StoredValue string
}

// Error describes the error
func (e ErrAlreadyExists) Error() string {
	return "value already exists"
}
