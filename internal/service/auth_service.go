package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sandbox/db/models"
	"sandbox/internal/config"
	"sandbox/internal/form"
	"sandbox/internal/lib/jwt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	TokenCookieName string
	Logger          *slog.Logger
	DB              *models.Queries
	Config          *config.Config
}

type LoginResult struct {
	User   *models.User
	Token  string
	Expiry int64
}

func (s *AuthService) Login(ctx context.Context, form *form.LoginForm) (*LoginResult, error) {
	user, err := s.DB.GetUserByEmail(ctx, form.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, expiry, err := jwt.NewExpiringToken(&jwt.ExpiringTokenArgs{
		UserId:        user.ID.String(),
		Email:         user.Email,
		Role:          string(user.Role),
		JwtSecret:     s.Config.Auth.JwtSecret,
		ExpiryMinutes: s.Config.Auth.JwtExpiryMinutes,
	})

	result := &LoginResult{
		User:   &user,
		Token:  token,
		Expiry: expiry,
	}

	return result, nil
}

func (s AuthService) SetAuthCookies(w http.ResponseWriter, result *LoginResult) {
	expiry := time.Now().Add(s.Config.Auth.JwtExpiryMinutes)
	http.SetCookie(w, &http.Cookie{
		Name:     s.TokenCookieName,
		Value:    result.Token,
		Expires:  expiry,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}

func (s AuthService) RemoveAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.TokenCookieName,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}

func (s AuthService) CreateAccount(ctx context.Context, form *form.RegisterForm) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("failed to hash user password", "error", err.Error())
		return errors.New("failed to register user")
	}

	err = s.DB.CreateUser(ctx, models.CreateUserParams{
		ID:       uuid.New(),
		Email:    form.Email,
		Role:     models.UserRoleADMIN,
		Password: string(passwordHash),
		Name:     pgtype.Text{Valid: false},
	})

	if err != nil {
		// TODO: handle email_unique constraint.
		s.Logger.Error("failed to save user to database", "error", err.Error())
		return errors.New("failed to register user")
	}

	return nil
}
