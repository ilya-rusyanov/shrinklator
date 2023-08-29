package storage

import (
	"fmt"
	"os"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

func NewFile(log *logger.Log, path string) (*os.File, error) {
	file, err := os.OpenFile(path,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0640)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	log.Info("opened file persistence",
		zap.String("file path", path))
	return file, nil
}
