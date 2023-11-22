package main

type Logger interface {
	Error(...any)
	Info(...any)
}

func printDeleteErrors(log Logger, ch <-chan error) {
	for err := range ch {
		log.Error("error while deleting: ", err.Error())
	}
}
