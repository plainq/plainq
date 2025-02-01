package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Metrics represents HTTP metrics collecting middlewares.
func Metrics() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx := chi.RouteContext(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			status := strconv.Itoa(ww.Status())
			route := ctx.RoutePattern()

			httpReqDur := httpReqDurationStr(r.Method, route, status)
			httpReqTotal := httpReqTotalStr(r.Method, route, status)

			metrics.GetOrCreateSummaryExt(httpReqDur, 5*time.Minute, []float64{0.95, 0.99}).
				UpdateDuration(start)

			metrics.GetOrCreateCounter(httpReqTotal).
				Inc()
		}

		return http.HandlerFunc(fn)
	}
}

func httpReqDurationStr(method, route, status string) string {
	return `http_request_duration{method="` + method + `", route="` + route + `", code="` + status + `"}`
}

func httpReqTotalStr(method, route, status string) string {
	return `http_requests_total{method="` + method + `", route="` + route + `", code="` + status + `"}`
}
