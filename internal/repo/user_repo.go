package repo

import (
	"app/internal/models"
	"context"
	"errors"

	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewUserRepo(db *sqlx.DB, logger *slog.Logger) *UserRepo {
	return &UserRepo{db, logger}
}

const createUserQuery string = `
	insert into "user" (id, email, role, password, name)
	values ($1, $2, $3, $4, $5)
`

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(ctx, createUserQuery, user.Id, user.Email, user.Role, user.Password, user.Name)
	if err != nil {
		e := err.(*pq.Error)
		r.logger.Error("failed to create user", "error", err.Error())

		switch e.Constraint {
		case models.ConstraintsUserEmailUnique:
			return errors.New("user with the provided email address already exists")

		default:
			return errors.New("failed to create user")
		}
	}

	return nil
}

const findUserByEmailQuery string = `
	select * from "user"
	where email = $1 and deleted_at is null
	limit 1
`

func (r *UserRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.QueryRowxContext(ctx, findUserByEmailQuery, email).StructScan(&user); err != nil {
		r.logger.Error("failed to find user by email", "error", err.Error())
		return nil, errors.New("failed to find user")
	}

	return &user, nil
}

const listUsersQuery string = `
	select *, COUNT(*) OVER() AS total_count from "user"
	where role = $1 and deleted_at is null
	order by created_at desc
	limit $2
	offset $3
`

type ListUsersResult struct {
	models.User
	TotalCount int `db:"total_count"`
}

type ListUsersArgs struct {
	Role   models.UserRole
	Limit  int
	Offset int
}

func (r *UserRepo) ListUsers(ctx context.Context, args *ListUsersArgs) ([]*ListUsersResult, error) {
	result := []*ListUsersResult{}
	if err := r.db.Select(&result, listUsersQuery, args.Role, args.Limit, args.Offset); err != nil {
		r.logger.Error("failed to list users", "error", err.Error())
		return nil, errors.New("failed to list users")
	}

	return result, nil
}
