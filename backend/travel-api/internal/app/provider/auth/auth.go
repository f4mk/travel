package auth

import (
	"context"

	authUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/auth"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Storer struct {
	repo *sqlx.DB
	log  *zerolog.Logger
}

func NewStorer(l *zerolog.Logger, r *sqlx.DB) *Storer {
	return &Storer{repo: r, log: l}
}

func (s *Storer) QueryByEmail(ctx context.Context, email string) (authUsecase.User, error) {
	u := authUsecase.User{}
	q := `SELECT * FROM users WHERE email = $1`
	if err := s.repo.GetContext(ctx, &u, q, email); err != nil {
		return authUsecase.User{}, err
	}
	return u, nil
}

func (s *Storer) QueryByID(ctx context.Context, id string) (authUsecase.User, error) {
	u := authUsecase.User{}
	q := `SELECT * FROM users WHERE user_id = $1`
	if err := s.repo.GetContext(ctx, &u, q, id); err != nil {
		return authUsecase.User{}, err
	}
	return u, nil
}

func (s *Storer) Update(ctx context.Context, u authUsecase.User) error {
	q := `UPDATE users SET
	name = :name, email = :email, is_active = :is_active,
	token_version = :token_version,
	roles = :roles, password_hash = :password_hash,
	date_updated = :date_updated
	WHERE user_id = :user_id;`
	_, err := s.repo.NamedExecContext(ctx, q, u)
	return err
}

func (s *Storer) StoreResetToken(ctx context.Context, rt authUsecase.ResetToken) error {
	q := `INSERT INTO reset_tokens (token_id, user_id, email, expires_at, issued_at)
	VALUES (:token_id, :user_id, :email, :expires_at, :issued_at)
	`
	_, err := s.repo.NamedExecContext(ctx, q, rt)
	return err
}

func (s *Storer) DeleteResetTokensByUserID(ctx context.Context, uID string) error {
	q := `DELETE from reset_tokens WHERE user_id = $1;`
	_, err := s.repo.ExecContext(ctx, q, uID)
	return err
}

func (s *Storer) QueryResetTokenByID(ctx context.Context, t string) (authUsecase.ResetToken, error) {
	rt := authUsecase.ResetToken{}
	q := `SELECT * FROM reset_tokens WHERE token_id = $1`
	if err := s.repo.GetContext(ctx, &rt, q, t); err != nil {
		return authUsecase.ResetToken{}, err
	}
	return rt, nil
}

func (s *Storer) DeleteToken(ctx context.Context, t authUsecase.DeleteToken) error {
	q := `INSERT INTO revoked_tokens (token_id, subject, token_version, issued_at, expires_at, revoked_at) 
	VALUES (:token_id, :subject, :token_version, :issued_at, :expires_at, :revoked_at)`
	_, err := s.repo.NamedExecContext(ctx, q, t)
	return err
}
