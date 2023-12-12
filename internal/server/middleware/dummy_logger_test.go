package middleware

type dummyLogger struct {
}

func (d *dummyLogger) Info(...any) {
}

func (d *dummyLogger) Error(...any) {
}

func (d *dummyLogger) Warnf(string, ...any) {
}
