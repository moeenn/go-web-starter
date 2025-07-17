package controller

import (
	"app/internal/form"
	"app/internal/lib/middleware"
	"app/internal/service"
	"app/views/components"
	"app/views/pages"
	"log/slog"
	"net/http"
)

type AuthController struct {
	Logger         *slog.Logger
	AuthService    *service.AuthService
	AuthMiddleware *middleware.AuthMiddleware
}

func NewAuthController(logger *slog.Logger, authService *service.AuthService, authMiddleware *middleware.AuthMiddleware) *AuthController {
	return &AuthController{
		Logger:         logger,
		AuthService:    authService,
		AuthMiddleware: authMiddleware,
	}
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/login", c.AuthMiddleware.NotLoggedIn(c.LoginPage))
	mux.HandleFunc("POST /auth/login", c.AuthMiddleware.NotLoggedIn(c.ProcessLoginRequest))
	mux.HandleFunc("GET /auth/register", c.AuthMiddleware.NotLoggedIn(c.RegisterPage))
	mux.HandleFunc("POST /auth/register", c.AuthMiddleware.NotLoggedIn(c.ProcessRegisterRequest))
	mux.HandleFunc("GET /auth/forgot-password", c.AuthMiddleware.NotLoggedIn(c.ForgotPasswordPage))
	mux.HandleFunc("POST /auth/forgot-password", c.AuthMiddleware.NotLoggedIn(c.ProcessForgotPasswordRequest))
	mux.HandleFunc("GET /auth/logout", c.AuthMiddleware.LoggedIn(c.ProcessLogoutRequest))
}

func (c AuthController) LoginPage(w http.ResponseWriter, r *http.Request) {
	html := pages.LoginPage(&pages.LoginPageProps{})
	html.Render(r.Context(), w)
}

func (c AuthController) ProcessLoginRequest(w http.ResponseWriter, r *http.Request) {
	form := form.LoginFormFromRequest(r)
	if err := form.Validate(); err != nil {
		html := components.LoginForm(components.LoginFormProps{
			Errors: err,
			Values: components.LoginFormFields{
				Email:    form.Email,
				Password: form.Password,
			},
		})

		html.Render(r.Context(), w)
		return
	}

	loginResult, err := c.AuthService.Login(r.Context(), &form)
	if err != nil {
		c.Logger.Error("failed to log in", "error", err.Error())

		message := components.Message{
			Message: "Invalid email or password",
			Type:    components.MessageTypeError,
		}
		html := components.LoginForm(components.LoginFormProps{
			Message: &message,
			Values: components.LoginFormFields{
				Email: form.Email,
			},
		})

		html.Render(r.Context(), w)
		return
	}

	c.AuthService.SetAuthCookies(w, loginResult)
	w.Header().Set("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusNoContent)
}

func (c AuthController) RegisterPage(w http.ResponseWriter, r *http.Request) {
	html := pages.RegisterPage()
	html.Render(r.Context(), w)
}

func (c AuthController) ProcessRegisterRequest(w http.ResponseWriter, r *http.Request) {
	form := form.RegisterFormFromRequest(r)
	if err := form.Validate(); err != nil {
		html := components.RegisterForm(components.RegisterFormProps{
			Errors: err,
			Values: components.RegisterFormFields{
				Email: form.Email,
			},
		})
		html.Render(r.Context(), w)
		return
	}

	if err := c.AuthService.CreateAccount(r.Context(), &form); err != nil {
		message := components.Message{Message: err.Error(), Type: components.MessageTypeError}
		html := components.RegisterForm(components.RegisterFormProps{
			Message: &message,
		})
		html.Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusNoContent)
}

func (c AuthController) ForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	html := pages.ForgotPasswordPage(pages.ForgotPasswordPageProps{})
	html.Render(r.Context(), w)
}

func (c AuthController) ProcessForgotPasswordRequest(w http.ResponseWriter, r *http.Request) {
	form := form.ForgotPasswordFormFromRequest(r)
	if err := form.Validate(); err != nil {
		html := components.ForgotPasswordForm(components.ForgotPasswordFormProps{
			Errors: err,
		})
		html.Render(r.Context(), w)
		return
	}

	message := components.Message{
		Message: "You request has been submitted. You will receive an email shortly with instructions to reset your password",
		Type:    components.MessageTypeSuccess,
	}
	html := components.ForgotPasswordForm(components.ForgotPasswordFormProps{
		Message: &message,
	})
	html.Render(r.Context(), w)
}

func (c AuthController) ProcessLogoutRequest(w http.ResponseWriter, r *http.Request) {
	c.AuthService.RemoveAuthCookies(w)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)
}
