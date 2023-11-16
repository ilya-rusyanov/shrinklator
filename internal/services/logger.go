package services

// Logger - logger for services
type Logger interface {
	Infof(string, ...any)
}
