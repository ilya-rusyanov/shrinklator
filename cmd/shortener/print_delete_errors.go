package main

// Logging object
type Logger interface {
	Error(...any)
	Info(...any)
	Warnf(string, ...any)
	Fatalf(string, ...any)
}

func printDeleteErrors(log Logger, ch <-chan error) {
	for err := range ch {
		log.Error("error while deleting: ", err.Error())
	}
}
