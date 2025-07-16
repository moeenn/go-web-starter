package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey struct{}

var jwtClaimsContextKey ContextKey = ContextKey{}

func JwtClaimsContextKey() ContextKey {
	return jwtClaimsContextKey
}

type JwtClaims struct {
	UserId string
	Email  string
	Role   string
}

func jwtClaimsFromMap(claims map[string]any) (*JwtClaims, error) {
	userId, userIdOk := claims["userId"].(string)
	if !userIdOk {
		return nil, errors.New("invalid or missing userId")
	}

	email, emailOk := claims["email"].(string)
	if !emailOk {
		return nil, errors.New("invalid or missing email")
	}

	role, roleOk := claims["role"].(string)
	if !roleOk {
		return nil, errors.New("invalid or missing role")
	}

	return &JwtClaims{
		UserId: userId,
		Email:  email,
		Role:   role,
	}, nil
}

func ValidateAndParseJwtClaims(jwtSecret []byte, bearerToken string) (*JwtClaims, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt claims")
	}

	// check token expiration.
	exp, err := claims.GetExpirationTime()
	if err != nil || exp.Time.Before(time.Now()) {
		return nil, errors.New("expired jwt")
	}

	// access claims.
	parsedClaims, err := jwtClaimsFromMap(claims)
	if err != nil {
		return nil, fmt.Errorf("invalid jwt claims: %w", err)
	}

	return parsedClaims, nil
}

func GetClaimsContext(baseCtx context.Context, jwtSecret []byte, token string) (context.Context, error) {
	claims, err := ValidateAndParseJwtClaims(jwtSecret, token)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(baseCtx, jwtClaimsContextKey, claims)
	return ctx, nil
}

func GetJwtClaims(r *http.Request) (*JwtClaims, bool) {
	claims := r.Context().Value(jwtClaimsContextKey)
	claimsCast, ok := claims.(*JwtClaims)
	if !ok {
		return nil, false
	}
	return claimsCast, true
}
