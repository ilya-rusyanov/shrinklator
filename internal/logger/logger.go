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

// Info logs with Info severity
func (l *Log) Info(args ...any) {
	l.logger.Sugar().Info(args)
}

// Infof logs with Info severity
func (l *Log) Infof(s string, args ...any) {
	l.logger.Sugar().Infof(s, args)
}

// Error logs with Error severity
func (l *Log) Error(args ...any) {
	l.logger.Sugar().Error(args)
}

// Debug logs with Debug severity
func (l *Log) Debug(args ...any) {
	l.logger.Sugar().Debug(args)
}

// Warn logs with Warn severity
func (l *Log) Warn(args ...any) {
	l.logger.Sugar().Warn(args)
}

// Warnf logs with Warn severity
func (l *Log) Warnf(s string, args ...any) {
	l.logger.Sugar().Warnf(s, args)
}

// Fatal logs with Fatal severity and shutdowns app
func (l *Log) Fatal(args ...any) {
	l.logger.Sugar().Fatal(args)
}

// Fatal logs with Fatal severity and shutdowns app
func (l *Log) Fatalf(s string, args ...any) {
	l.logger.Sugar().Fatal(args)
}
