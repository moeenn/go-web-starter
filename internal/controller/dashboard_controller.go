package controller

import (
	"log/slog"
	"net/http"
	"sandbox/views/pages"
)

type DashboardController struct {
	Logger *slog.Logger
}

func NewDashboardController(logger *slog.Logger) *DashboardController {
	return &DashboardController{
		Logger: logger,
	}
}

func (c *DashboardController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/dashboard", c.DashboardHomePage)
	mux.HandleFunc("/dashboard/users", c.DashboardUsersPage)
	mux.HandleFunc("/dashboard/clients", c.DashboardClientsPage)
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
