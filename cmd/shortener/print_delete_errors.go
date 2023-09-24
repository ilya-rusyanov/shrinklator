package main

import (
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

func printDeleteErrors(log *logger.Log, ch <-chan error) {
	for err := range ch {
		log.Error("error while deleting", zap.String("message", err.Error()))
	}
}
