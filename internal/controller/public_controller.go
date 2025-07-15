package controller

import (
	"log/slog"
	"sandbox/internal/lib"
	"sandbox/views/pages"

	"github.com/labstack/echo/v4"
)

type PublicController struct {
	Logger *slog.Logger
}

func NewPublicController(logger *slog.Logger) *PublicController {
	return &PublicController{
		Logger: logger,
	}
}

func (c *PublicController) RegisterRoutes(e *echo.Echo) {
	e.GET("/", c.HomePage)
}

func (c PublicController) HomePage(ctx echo.Context) error {
	html := pages.HomePage(pages.HomePageProps{
		IsLoggedIn:     false,
		FlashedMessage: lib.Ref("Welcome to our simple home-page"),
	})

	return html.Render(ctx.Request().Context(), ctx.Response().Writer)
}
