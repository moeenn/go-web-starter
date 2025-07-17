package repo

import (
	"context"
	"sandbox/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db}
}

const createUserQuery string = `
	insert into "user" (id, email, role, password, name)
	values ($1, $2, $3, $4, $5)
`

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(ctx, createUserQuery, user.Id, user.Email, user.Role, user.Password, user.Name)
	// TODO: handle constraint violations.
	return err
}

const findUserByEmailQuery string = `
	select * from "user"
	where email = $1 and deleted_at is null
	limit 1
`

func (r *UserRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.QueryRowxContext(ctx, findUserByEmailQuery, email).StructScan(&user); err != nil {
		// TODO: wrap db error.
		return nil, err
	}

	return &user, nil
}
