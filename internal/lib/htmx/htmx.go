package htmx

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func IsHTMXRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

func Redirect(c echo.Context, url string) error {
	c.Response().Header().Set("HX-Location", url)
	return c.NoContent(http.StatusNoContent)
}
