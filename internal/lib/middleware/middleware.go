package middleware

import (
	"net/http"
	"sandbox/internal/lib/htmx"

	"github.com/labstack/echo/v4"
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

func (m AuthMiddleware) LoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		redirectUrl := "/"
		_, err := c.Cookie(m.TokenCookieName)
		if err != nil {
			if htmx.IsHTMXRequest(c) {
				return htmx.Redirect(c, redirectUrl)
			}
			return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		}

		return next(c)
	}
}

func (m AuthMiddleware) NotLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		redirectUrl := "/dashboard"
		if _, err := c.Cookie(m.TokenCookieName); err == nil {
			if htmx.IsHTMXRequest(c) {
				c.Response().Header().Set("HX-Location", redirectUrl)
				return c.NoContent(http.StatusNoContent)
			}

			return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		}

		return next(c)
	}
}

// TODO: fixme.
// func (m AuthMiddleware) SetJwtClaimsInContext(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		cookie, err := c.Cookie(m.TokenCookieName)
// 		fmt.Printf("\n\n cookie: %v, err: %v\n\n", cookie, err)

// 		if err == nil && cookie != nil {
// 			claims, err := validateAndParseJwtClaims(m.JwtSecret, cookie.Value)
// 			if err == nil {
// 				c.Set(jwtClaimsContextKey, claims)
// 			}
// 		}

// 		return next(c)
// 	}
// }
