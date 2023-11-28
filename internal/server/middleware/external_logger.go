package middleware

// ExternalLogger - logger that will be used by middleware to log messages
type ExternalLogger interface {
	Info(...any)
	Error(...any)
	Warnf(string, ...any)
}
