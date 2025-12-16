package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JasonLeonnn/jalytics/internal/metrics"
	"github.com/go-chi/chi/v5"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		route := r.URL.Path
		if rc := chi.RouteContext(r.Context()); rc != nil {
			if pat := rc.RoutePattern(); pat != "" {
				route = pat
			}
		}

		metrics.HttpRequestDuration.WithLabelValues(
			route,
			r.Method,
			fmt.Sprintf("%d", rr.statusCode),
		).Observe(time.Since(start).Seconds())
	})
}
