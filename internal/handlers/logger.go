package handlers

// Logger - logger for handlers
type Logger interface {
	Info(...any)
	Infof(string, ...any)
	Error(...any)
}
