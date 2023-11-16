package storage

import (
	"context"

	"github.com/ilya-rusyanov/shrinklator/internal/config"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
)

// Factory - application storage factory
type Factory struct {
}

// MustInitStorage creates storage or panics
func (f *Factory) MustInitStorage(config config.Config, log *logger.Log) Interface {
	switch {
	case config.StoreInDB:
		ctx := context.Background()
		db, err := NewPostgres(ctx, log, config.DSN)
		if err != nil {
			panic(err)
		}

		log.Info("storage is database")
		return db
	case config.StoreInFile:
		file, err := NewFile(log, config.FileStoragePath)
		if err != nil {
			panic(err)
		}

		log.Info("storage is file")
		return file
	default:
		memory := NewInMemory(log)
		log.Info("in memory storage")
		return memory
	}
}
