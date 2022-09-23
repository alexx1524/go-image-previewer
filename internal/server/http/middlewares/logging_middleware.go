package middlewares

import (
	"net/http"
	"time"

	"github.com/alexx1524/go-image-previewer/internal/logger"
)

type LoggingMiddleware struct {
	Logger logger.Logger
}

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lwr *loggingWriter) WriteHeader(code int) {
	lwr.statusCode = code
	lwr.ResponseWriter.WriteHeader(code)
}

func (logMiddleware *LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		logWriter := &loggingWriter{writer, http.StatusOK}
		next.ServeHTTP(writer, request)

		logMiddleware.Logger.LogHTTPRequest(request, logWriter.statusCode, time.Since(start))
	})
}
