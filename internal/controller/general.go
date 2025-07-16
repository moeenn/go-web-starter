package controller

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, html templ.Component) error {
	return html.Render(c.Request().Context(), c.Response().Writer)
}
