package middleware

import (
	"net/http"
	"time"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

// Logger - middleware for logging HTTP requests
type Logger struct {
	log *logger.Log
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

// NewLogger constructs Logger objects
func NewLogger(log *logger.Log) *Logger {
	return &Logger{
		log: log,
	}
}

// Middleware creates actual middleware
func (l *Logger) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}
			lw := loggingResponseWriter{
				ResponseWriter: rw, // встраиваем оригинальный http.ResponseWriter
				responseData:   responseData,
			}
			next.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

			duration := time.Since(start)

			contentType := r.Header.Get("Content-Type")
			acceptEncoding := r.Header.Get("Accept-Encoding")
			contentEncoding := r.Header.Get("Content-Encoding")

			l.log.Info(
				"recieved request",
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.Int("status", responseData.status), // получаем перехваченный код статуса ответа
				zap.Int64("duration", int64(duration)),
				zap.Int("size", responseData.size), // получаем перехваченный размер ответа
				zap.String("content-type", contentType),
				zap.String("content-encoding", contentEncoding),
				zap.String("accept-encoding", acceptEncoding),
			)
		})
	}
}

// Write writes original response
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

// WriteHeader writes header
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}
