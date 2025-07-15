package controller

import (
	"log/slog"
	"sandbox/views/pages"

	"github.com/labstack/echo/v4"
)

type DashboardController struct {
	Logger *slog.Logger
}

func NewDashboardController(logger *slog.Logger) *DashboardController {
	return &DashboardController{
		Logger: logger,
	}
}

func (c *DashboardController) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/dashboard")
	g.GET("", c.DashboardHomePage)
	g.GET("/users", c.DashboardUsersPage)
	g.GET("/clients", c.DashboardClientsPage)
}

func (c DashboardController) DashboardHomePage(ctx echo.Context) error {
	html := pages.DashboardHomePage(pages.DashboardHomePageProps{
		CurrentUrl: ctx.Path(),
	})

	return html.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func (c DashboardController) DashboardUsersPage(ctx echo.Context) error {
	html := pages.DashboardUsersPage(pages.DashboardUsersPageProps{
		CurrentUrl: ctx.Path(),
	})

	return html.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func (c DashboardController) DashboardClientsPage(ctx echo.Context) error {
	html := pages.DashboardClientsPage(pages.DashboardClientsPageProps{
		CurrentUrl: ctx.Path(),
	})

	return html.Render(ctx.Request().Context(), ctx.Response().Writer)
}
