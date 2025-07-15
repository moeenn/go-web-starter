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
		if _, err := c.Cookie(m.TokenCookieName); err != nil {
			c.Response().Header().Set("HX-Redirect", "/auth/login")
			return c.NoContent(http.StatusNoContent)
		}

		return next(c)
	}
}

func (m AuthMiddleware) NotLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, err := c.Cookie(m.TokenCookieName); err == nil {
			c.Response().Header().Set("HX-Redirect", "/dashboard")
			return c.NoContent(http.StatusNoContent)
		}

		return next(c)
	}
}
