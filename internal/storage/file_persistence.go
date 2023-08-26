package storage

import (
	"fmt"
	"os"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

type FilePersistence struct {
	file *os.File
}

func NewFilePersistence(filePath string) (*FilePersistence, error) {
	file, err := os.OpenFile(filePath,
		os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_SYNC,
		0640)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	res := &FilePersistence{file: file}
	logger.Log.Info("created file persistence",
		zap.String("file path", filePath))
	return res, nil
}

func (p *FilePersistence) Append(string, string) {
}

func (p *FilePersistence) ReadAll() (values map[string]string, err error) {
	values = make(map[string]string)
	return
}
