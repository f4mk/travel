package auth

import (
	"context"

	authUsecase "github.com/f4mk/api/internal/app/usecase/auth"

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

func (r *Repo) DeleteToken(_ context.Context, t authUsecase.DeleteToken) error {
	query := `INSERT INTO revoked_tokens (token_id, subject, issued_at, expires_at, revoked_at) 
	VALUES (:token_id, :subject, :issued_at, :expires_at, :revoked_at)`
	_, err := r.repo.NamedExec(query, t)
	return err
}

func (r *Repo) QueryByEmail(_ context.Context, email string) (authUsecase.User, error) {
	u := authUsecase.User{}
	q := `SELECT * FROM users WHERE email = $1`
	if err := r.repo.Get(&u, q, email); err != nil {
		return authUsecase.User{}, err
	}
	return u, nil
}

func (r *Repo) QueryByID(_ context.Context, id string) (authUsecase.User, error) {
	u := authUsecase.User{}
	q := `SELECT * FROM users WHERE user_id = $1`
	if err := r.repo.Get(&u, q, id); err != nil {
		return authUsecase.User{}, err
	}
	return u, nil
}

func (r *Repo) Update(ctx context.Context, u authUsecase.User) error {
	q := `UPDATE users SET name = :name, email = :email, roles = :roles, password_hash = :password_hash, date_updated = :date_updated
			WHERE user_id = :user_id;`
	_, err := r.repo.NamedExecContext(ctx, q, u)
	if err != nil {
		return err
	}
	return nil
}

// TODO: implement
//
//revive:disable
func (r *Repo) DeleteAllTokes(_ context.Context, uID string) error {
	return nil
}

//revive:enable
