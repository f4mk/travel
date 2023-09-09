package user

import (
	"context"

	"github.com/f4mk/travel/backend/travel-api/internal/app/usecase/user"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

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
	q := `INSERT INTO users(user_id, name, email, is_active, token_version, roles, password_hash, date_created, date_updated) 
			VALUES(:user_id, :name, :email, :is_active, :token_version, :roles, :password_hash, :date_created, :date_updated);`
	_, err := r.repo.NamedExecContext(ctx, q, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) QueryByID(ctx context.Context, userID string) (user.User, error) {
	res := user.User{}
	q := `SELECT * from users WHERE user_id = $1`
	if err := r.repo.GetContext(ctx, &res, q, userID); err != nil {
		return user.User{}, err
	}
	return res, nil
}

func (r *Repo) Update(ctx context.Context, u user.User) error {
	q := `UPDATE users SET 
					name = :name, email = :email, 
					roles = :roles, password_hash = :password_hash, 
					date_updated = :date_updated
				WHERE user_id = :user_id;`
	_, err := r.repo.NamedExecContext(ctx, q, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, user user.User) error {
	q := `UPDATE users SET 
					is_active = :is_active, 
					date_updated = :date_updated,
					token_version = :token_version
				WHERE user_id = :user_id`
	_, err := r.repo.NamedExecContext(ctx, q, user)
	if err != nil {
		return err
	}
	return nil
}
