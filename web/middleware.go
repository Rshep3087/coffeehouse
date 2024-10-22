package web

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// loggingmw is a middleware that logs the request and response
// of an HTTP request.
func loggingmw(log *zap.SugaredLogger, h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		logw := log.With(
			"method", r.Method,
			"uri", r.RequestURI,
		)

		logw.Info("request started")

		h.ServeHTTP(w, r)

		duration := time.Since(t).Milliseconds()
		logw.Infow("request completed", "duration", duration)
	})
}
