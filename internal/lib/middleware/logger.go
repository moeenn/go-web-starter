package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func Logger(logger *slog.Logger, next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// don't log static asset requests.
		if strings.HasPrefix(path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)

		method := r.Method
		elapsed := time.Since(start)

		logger.Info(
			"request",
			"method", method,
			"path", path,
			"elapsed", elapsed.Milliseconds(),
		)
	}

	return http.HandlerFunc(f)
}
