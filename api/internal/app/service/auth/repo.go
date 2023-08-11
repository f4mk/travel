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

func NewRepo(r *sqlx.DB, l *zerolog.Logger) *Repo {

	return &Repo{repo: r, log: l}
}

func (r *Repo) DeleteToken(ctx context.Context, t authUsecase.DeleteToken) error {
	query := `INSERT INTO revoked_tokens (token_id, subject, issued_at, expires_at, revoked_at) 
	VALUES (:token_id, :subject, :issued_at, :expires_at, :revoked_at)`
	_, err := r.repo.NamedExec(query, t)
	return err
}

// TODO: implement
func (r *Repo) DeleteAllTokes(ctx context.Context, uID string) error {
	return nil
}

func (r *Repo) QueryByEmail(ctx context.Context, email string) (authUsecase.User, error) {
	user := authUsecase.User{}
	q := `SELECT * FROM users WHERE email = $1`
	if err := r.repo.Get(&user, q, email); err != nil {
		return authUsecase.User{}, err
	}
	return user, nil
}

func (r *Repo) QueryByID(ctx context.Context, id string) (authUsecase.User, error) {
	user := authUsecase.User{}
	q := `SELECT * FROM users WHERE user_id = $1`
	if err := r.repo.Get(&user, q, id); err != nil {
		return authUsecase.User{}, err
	}
	return user, nil
}
