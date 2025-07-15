package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	TokenCookieName string
}

func NewAuthMiddleware(cookieName string) *AuthMiddleware {
	return &AuthMiddleware{
		TokenCookieName: cookieName,
	}
}

func (m AuthMiddleware) LoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		redirectUrl := "/"

		if _, err := c.Cookie(m.TokenCookieName); err != nil {
			if IsHTMXRequest(c) {
				c.Response().Header().Set("HX-Location", redirectUrl)
				return c.NoContent(http.StatusNoContent)
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
			if IsHTMXRequest(c) {
				c.Response().Header().Set("HX-Location", redirectUrl)
				return c.NoContent(http.StatusNoContent)
			}

			return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		}

		return next(c)
	}
}

func IsHTMXRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}
