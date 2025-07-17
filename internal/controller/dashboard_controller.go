package controller

import (
	"log/slog"
	"net/http"
	"sandbox/internal/lib/middleware"
	"sandbox/internal/models"
	"sandbox/internal/repo"
	"sandbox/views/pages"
	"strconv"
)

type DashboardController struct {
	Logger         *slog.Logger
	UserRepo       *repo.UserRepo
	AuthMiddleware *middleware.AuthMiddleware
}

func NewDashboardController(logger *slog.Logger, authMiddleware *middleware.AuthMiddleware, userRepo *repo.UserRepo) *DashboardController {
	return &DashboardController{
		Logger:         logger,
		AuthMiddleware: authMiddleware,
		UserRepo:       userRepo,
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
	limit, offset := GetLimitAndOffset(r)

	result, err := c.UserRepo.ListUsers(r.Context(), &repo.ListUsersArgs{
		Role:   models.UserRoleClient,
		Limit:  limit,
		Offset: offset,
	})

	// TODO: improve.
	if err != nil {
		c.Logger.Error("failed to list users", "error", err.Error())
		w.Write([]byte(`<h1>Something went wrong</h1>`))
		return
	}

	totalCount := 0
	if len(result) > 0 {
		totalCount = result[0].TotalCount
	}

	clients := make([]*models.User, len(result))
	for i := range len(result) {
		clients[i] = &result[i].User
	}

	url := r.URL.Path
	html := pages.DashboardClientsPage(pages.DashboardClientsPageProps{
		CurrentUrl: url,
		TableData: &pages.ClientsTableProps{
			Limit:      limit,
			Offset:     offset,
			TotalCount: totalCount,
			Url:        url,
			Clients:    clients,
		},
	})

	html.Render(r.Context(), w)
}

func GetLimitAndOffset(r *http.Request) (int, int) {
	limit := 10
	offset := 0

	q := r.URL.Query()
	limitRaw := q.Get("limit")
	if limitRaw != "" {
		parsedLimit, err := strconv.Atoi(limitRaw)
		if err == nil && parsedLimit >= 1 {
			limit = parsedLimit
		}
	}

	offsetRaw := q.Get("offset")
	if offsetRaw != "" {
		parsedoffset, err := strconv.Atoi(offsetRaw)
		if err == nil && parsedoffset >= 0 {
			offset = parsedoffset
		}
	}

	return limit, offset
}
