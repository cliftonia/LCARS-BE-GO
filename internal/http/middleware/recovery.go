package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/cliftonbaggerman/subspace-backend/internal/logger"
)

// Recovery recovers from panics and returns a 500 error
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log := logger.Get()
				log.Error("Panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
					"path", r.URL.Path,
					"method", r.Method,
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
