package logrus

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Log - application logger
type Log struct {
	logger *logrus.Logger
}

// NewLogger constructs Log object
func NewLogger(level string) (*Log, error) {
	res := Log{
		logger: logrus.New(),
	}

	res.logger.SetFormatter(&logrus.JSONFormatter{})

	l, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse logging level: %w", err)
	}

	res.logger.SetLevel(l)

	return &res, nil
}

// Info logs with Info severity
func (l *Log) Info(args ...any) {
	l.logger.Info(args...)
}

// Infof logs with Info severity
func (l *Log) Infof(s string, args ...any) {
	l.logger.Infof(s, args...)
}

// Error logs with Error severity
func (l *Log) Error(args ...any) {
	l.logger.Error(args...)
}

// Debug logs with Debug severity
func (l *Log) Debug(args ...any) {
	l.logger.Debug(args...)
}

// Warn logs with Warn severity
func (l *Log) Warn(args ...any) {
	l.logger.Warn(args...)
}
