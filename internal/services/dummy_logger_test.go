package services

type dummyLogger struct {
}

func (l *dummyLogger) Debug(...any) {
}

func (l *dummyLogger) Info(...any) {
}

func (l *dummyLogger) Infof(string, ...any) {
}

func (l *dummyLogger) Warn(...any) {
}

func (l *dummyLogger) Warnf(string, ...any) {
}

func (l *dummyLogger) Error(...any) {
}
