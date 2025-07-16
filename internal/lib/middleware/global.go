package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"sandbox/internal/lib/jwt"
	"time"
)

func (m AuthMiddleware) SetClaimsContext(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenCookie, err := r.Cookie(m.TokenCookieName)
		if err == nil {
			claims, err := jwt.ValidateAndParseJwtClaims(m.JwtSecret, tokenCookie.Value)
			if err == nil {
				ctx = context.WithValue(ctx, jwt.JwtClaimsContextKey(), claims)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(f)
}

func Logger(logger *slog.Logger, next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		path := r.URL.Path
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
