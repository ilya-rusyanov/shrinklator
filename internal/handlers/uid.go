package handlers

// ContextKey is a type for handling values through Context
type ContextKey int

// available keys
const (
	// UID is for user IDs
	UID ContextKey = iota
)
