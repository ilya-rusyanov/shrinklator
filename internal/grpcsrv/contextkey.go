package grpcsrv

type contextKey int

// available keys
const (
	// UID is for user IDs
	uid contextKey = iota
)
