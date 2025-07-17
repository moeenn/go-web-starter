package form

import (
	"app/views/components"
	"net/http"
	"strings"
)

type LoginForm struct {
	Email    string
	Password string
}

func LoginFormFromRequest(r *http.Request) LoginForm {
	email := r.FormValue("email")
	password := r.FormValue("password")

	return LoginForm{
		Email:    email,
		Password: password,
	}
}

func (f LoginForm) Validate() *components.LoginFormFields {
	errs := &components.LoginFormFields{}
	isValid := true

	if f.Email == "" || !strings.Contains(f.Email, "@") {
		errs.Email = "Please provide a valid email address"
		isValid = false
	}

	if f.Password == "" || len(f.Password) < 8 {
		errs.Password = "Please provide a stonger password"
		isValid = false
	}

	if !isValid {
		return errs
	}

	return nil
}

type ForgotPasswordForm struct {
	Email string
}

func ForgotPasswordFormFromRequest(r *http.Request) ForgotPasswordForm {
	email := r.FormValue("email")
	return ForgotPasswordForm{
		Email: email,
	}
}

func (f ForgotPasswordForm) Validate() *components.ForgotPasswordFormFields {
	if f.Email == "" || !strings.Contains(f.Email, "@") {
		return &components.ForgotPasswordFormFields{
			Email: "Please provide a valid email address",
		}
	}

	return nil
}

type RegisterForm struct {
	Email           string
	Password        string
	ConfirmPassword string
}

func RegisterFormFromRequest(r *http.Request) RegisterForm {
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	return RegisterForm{
		Email:           email,
		Password:        password,
		ConfirmPassword: confirmPassword,
	}
}

func (f RegisterForm) Validate() *components.RegisterFormFields {
	errs := &components.RegisterFormFields{}
	isValid := true

	if f.Email == "" || !strings.Contains(f.Email, "@") {
		errs.Email = "Please provide a valid email address"
		isValid = false
	}

	if f.Password == "" || len(f.Password) < 8 {
		errs.Password = "Please provide a stonger password"
		isValid = false
	}

	if f.ConfirmPassword == "" || f.ConfirmPassword != f.Password {
		errs.ConfirmPassword = "Password confirmation failed"
		isValid = false
	}

	if !isValid {
		return errs
	}

	return nil
}
