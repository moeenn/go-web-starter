package controller

import (
	"log/slog"
	"net/http"
	"sandbox/internal/lib/middleware"
	"sandbox/views/pages"
)

type DashboardController struct {
	Logger         *slog.Logger
	AuthMiddleware *middleware.AuthMiddleware
}

func NewDashboardController(logger *slog.Logger, authMiddleware *middleware.AuthMiddleware) *DashboardController {
	return &DashboardController{
		Logger:         logger,
		AuthMiddleware: authMiddleware,
	}
}

func (c *DashboardController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/dashboard", c.AuthMiddleware.LoggedIn(c.DashboardHomePage))
	mux.HandleFunc("/dashboard/users", c.AuthMiddleware.LoggedIn(c.DashboardUsersPage))
	mux.HandleFunc("/dashboard/clients", c.AuthMiddleware.LoggedIn(c.DashboardClientsPage))
}

func (c DashboardController) DashboardHomePage(w http.ResponseWriter, r *http.Request) {
	html := pages.DashboardHomePage(pages.DashboardHomePageProps{
		CurrentUrl: r.URL.Path,
	})

	// TODO: handle render error.
	html.Render(r.Context(), w)
}

func (c DashboardController) DashboardUsersPage(w http.ResponseWriter, r *http.Request) {
	html := pages.DashboardUsersPage(pages.DashboardUsersPageProps{
		CurrentUrl: r.URL.Path,
	})

	html.Render(r.Context(), w)
}

func (c DashboardController) DashboardClientsPage(w http.ResponseWriter, r *http.Request) {
	html := pages.DashboardClientsPage(pages.DashboardClientsPageProps{
		CurrentUrl: r.URL.Path,
	})

	html.Render(r.Context(), w)
}
