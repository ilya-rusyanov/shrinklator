package handlers

// Logger - logger for handlers
type Logger interface {
	Debug(...any)
	Info(...any)
	Infof(string, ...any)
	Error(...any)
}
