package logger

import (
	"go.uber.org/zap"
)

// Log - application logger
type Log struct {
	logger *zap.Logger
}

// NewLogger constructs Log object
func NewLogger(level string) (*Log, error) {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	res := Log{
		logger: zl,
	}

	return &res, nil
}

func (l *Log) Info(args ...any) {
	l.logger.Sugar().Info(args)
}

func (l *Log) Infof(s string, args ...any) {
	l.logger.Sugar().Infof(s, args)
}

func (l *Log) Error(args ...any) {
	l.logger.Sugar().Error(args)
}

func (l *Log) Debug(args ...any) {
	l.logger.Sugar().Debug(args)
}

func (l *Log) Warn(args ...any) {
	l.logger.Sugar().Warn(args)
}
