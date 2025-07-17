package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleClient   UserRole = "CLIENT"
	UserRoleEmployee UserRole = "EMPLOYEE"
)

type User struct {
	Id        uuid.UUID      `db:"id"`
	Email     string         `db:"email"`
	Role      UserRole       `db:"role"`
	Password  string         `db:"password"`
	Name      sql.NullString `db:"name"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
}

type ConstraintsUser string

const (
	ConstraintsUserEmailUnique ConstraintsUser = "email_unique"
)
