package middleware

import (
	"context"
	"net/http"
	"sandbox/internal/lib/jwt"
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
