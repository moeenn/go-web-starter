package middleware

import (
	"net/http"
)

type AuthMiddleware struct {
	TokenCookieName string
	JwtSecret       []byte
}

func NewAuthMiddleware(cookieName string, jwtSecret []byte) *AuthMiddleware {
	return &AuthMiddleware{
		TokenCookieName: cookieName,
		JwtSecret:       jwtSecret,
	}
}

func (m AuthMiddleware) LoggedIn(next http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		redirectUrl := "/"
		_, err := r.Cookie(m.TokenCookieName)
		if err != nil {
			if IsHTMXRequest(r) {
				r.Header.Set("HX-Redirect", redirectUrl)
				w.WriteHeader(http.StatusNoContent)
				return
			}

			http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
			return
		}

		next(w, r)
	}

	return http.HandlerFunc(f)
}

func (m AuthMiddleware) NotLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		redirectUrl := "/dashboard"
		if _, err := r.Cookie(m.TokenCookieName); err == nil {
			if IsHTMXRequest(r) {
				w.Header().Set("HX-Redirect", redirectUrl)
				w.WriteHeader(http.StatusNoContent)
				return
			}

			http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
			return
		}

		next(w, r)
	}

	return http.HandlerFunc(f)
}

func IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
