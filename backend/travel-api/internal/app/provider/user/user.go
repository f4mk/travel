package user

import (
	"context"

	"github.com/f4mk/travel/backend/travel-api/internal/app/usecase/user"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// Storer is a type that's responsible for direct interaction with a DB.
// Only Storer should be aware of the data layout in a DB.
// No other logic but DB requests should be presented in this file.
// Storer methods should accept parameters that are mandatory for DB queries
// and return appropriate results

type Repo struct {
	repo *sqlx.DB
	log  *zerolog.Logger
}

func NewRepo(l *zerolog.Logger, r *sqlx.DB) *Repo {
	return &Repo{repo: r, log: l}
}

func (r *Repo) QueryAll(ctx context.Context) ([]user.User, error) {
	res := []user.User{}
	q := `SELECT * from users`
	if err := r.repo.SelectContext(ctx, &res, q); err != nil {
		return []user.User{}, err
	}
	return res, nil
}

func (r *Repo) Create(ctx context.Context, u user.User) error {
	q := `INSERT INTO users(user_id, name, email, token_version, roles, password_hash, date_created, date_updated) 
			VALUES(:user_id, :name, :email, :token_version, :roles, :password_hash, :date_created, :date_updated);`
	_, err := r.repo.NamedExecContext(ctx, q, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) QueryByID(ctx context.Context, uID string) (user.User, error) {
	res := user.User{}
	q := `SELECT * from users WHERE user_id = $1`
	if err := r.repo.GetContext(ctx, &res, q, uID); err != nil {
		return user.User{}, err
	}
	return res, nil
}

func (r *Repo) Update(ctx context.Context, u user.User) error {
	q := `UPDATE users SET name = :name, email = :email, roles = :roles, password_hash = :password_hash, date_updated = :date_updated
			WHERE user_id = :user_id;`
	_, err := r.repo.NamedExecContext(ctx, q, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, uID string) error {
	q := `DELETE from users WHERE user_id = $1;`
	_, err := r.repo.ExecContext(ctx, q, uID)
	if err != nil {
		return err
	}
	return nil
}
