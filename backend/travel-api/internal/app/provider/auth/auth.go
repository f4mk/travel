package auth

import (
	"context"

	authUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/auth"

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

func (r *Repo) QueryByEmail(ctx context.Context, email string) (authUsecase.User, error) {
	u := authUsecase.User{}
	q := `SELECT * FROM users WHERE email = $1`
	if err := r.repo.GetContext(ctx, &u, q, email); err != nil {
		return authUsecase.User{}, err
	}
	return u, nil
}

func (r *Repo) QueryByID(ctx context.Context, id string) (authUsecase.User, error) {
	u := authUsecase.User{}
	q := `SELECT * FROM users WHERE user_id = $1`
	if err := r.repo.GetContext(ctx, &u, q, id); err != nil {
		return authUsecase.User{}, err
	}
	return u, nil
}

func (r *Repo) Update(ctx context.Context, u authUsecase.User) error {
	q := `UPDATE users SET
	name = :name, email = :email, is_active = :is_active,
	token_version = :token_version,
	roles = :roles, password_hash = :password_hash,
	date_updated = :date_updated
	WHERE user_id = :user_id;`
	_, err := r.repo.NamedExecContext(ctx, q, u)
	return err
}

func (r *Repo) StoreResetToken(ctx context.Context, rt authUsecase.ResetToken) error {
	q := `INSERT INTO reset_tokens (token_id, user_id, email, expires_at, issued_at)
	VALUES (:token_id, :user_id, :email, :expires_at, :issued_at)
	`
	_, err := r.repo.NamedExecContext(ctx, q, rt)
	return err
}

func (r *Repo) DeleteResetTokensByUserID(ctx context.Context, uID string) error {
	q := `DELETE from reset_tokens WHERE user_id = $1;`
	_, err := r.repo.ExecContext(ctx, q, uID)
	return err
}

func (r *Repo) QueryResetTokenByID(ctx context.Context, t string) (authUsecase.ResetToken, error) {
	rt := authUsecase.ResetToken{}
	q := `SELECT * FROM reset_tokens WHERE token_id = $1`
	if err := r.repo.GetContext(ctx, &rt, q, t); err != nil {
		return authUsecase.ResetToken{}, err
	}
	return rt, nil
}

func (r *Repo) DeleteToken(ctx context.Context, t authUsecase.DeleteToken) error {
	q := `INSERT INTO revoked_tokens (token_id, subject, token_version, issued_at, expires_at, revoked_at) 
	VALUES (:token_id, :subject, :token_version, :issued_at, :expires_at, :revoked_at)`
	_, err := r.repo.NamedExecContext(ctx, q, t)
	return err
}
