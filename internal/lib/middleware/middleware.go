package middleware

import (
	"app/internal/lib/jwt"
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
		cookie, err := r.Cookie(m.TokenCookieName)
		if err != nil {
			redirect(w, r, redirectUrl)
			return
		}

		ctx, err := jwt.GetClaimsContext(r.Context(), m.JwtSecret, cookie.Value)
		if err != nil {
			redirect(w, r, redirectUrl)
			return
		}

		next(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(f)
}

func (m AuthMiddleware) NotLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie(m.TokenCookieName); err == nil {
			redirect(w, r, "/dashboard")
			return
		}
		next(w, r)
	}

	return http.HandlerFunc(f)
}

func redirect(w http.ResponseWriter, r *http.Request, redirectUrl string) {
	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", redirectUrl)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
