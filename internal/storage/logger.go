package storage

// Logger - logger for storage
type Logger interface {
	Info(...any)
	Warn(...any)
	Debug(...any)
}
