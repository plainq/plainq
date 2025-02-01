package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Recovery() func(next http.Handler) http.Handler {
	return middleware.Recoverer
}

func RedirectSlashes() func(next http.Handler) http.Handler {
	return middleware.RedirectSlashes
}

func Profiler() http.Handler {
	return middleware.Profiler()
}

func Cors(options cors.Options) func(next http.Handler) http.Handler {
	return cors.Handler(options)
}
