package middleware

import (
	"context"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	TraceID                      = "trace_id"
	ContextKeyTraceID ContextKey = "trace_id"
)

// RequestLogger HTTP middleware setting a value on the request context
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()

		// функция Now() возвращает текущее время
		start := time.Now()
		uri := r.RequestURI
		// метод запроса
		method := r.Method

		logger.Log.Info("REQUEST",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.String(TraceID, traceID),
		)

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		ctx := context.WithValue(r.Context(), ContextKeyTraceID, traceID)
		r = r.WithContext(ctx)

		next.ServeHTTP(&lw, r)

		duration := time.Since(start).String()

		logger.Log.Info("RESPONSE",
			zap.Int("status", responseData.status), // получаем перехваченный код статуса ответа
			zap.Int("size", responseData.size),     // получаем перехваченный размер ответа
			zap.String("duration", duration),
			zap.String(TraceID, traceID),
		)
	})
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

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	//nolint:wrapcheck // need to return err
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}
