package auth

import (
	"context"

	authUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"

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
	ctx, span := web.AddSpan(ctx, "provider.auth.query-by-email")
	defer span.End()
	user := StorerUser{}
	q := `SELECT * FROM users WHERE email = $1;`
	if err := s.repo.GetContext(ctx, &user, q, email); err != nil {
		return authUsecase.User{}, err
	}
	res := authUsecase.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		IsActive:     user.IsActive,
		IsDeleted:    user.IsDeleted,
		TokenVersion: user.TokenVersion,
		Roles:        user.Roles,
		PasswordHash: user.PasswordHash,
		DateCreated:  user.DateCreated,
		DateUpdated:  user.DateUpdated,
	}
	return res, nil
}

func (s *Storer) QueryByID(ctx context.Context, id string) (authUsecase.User, error) {
	ctx, span := web.AddSpan(ctx, "provider.auth.query-by-id")
	defer span.End()
	user := StorerUser{}
	q := `SELECT * FROM users WHERE user_id = $1;`
	if err := s.repo.GetContext(ctx, &user, q, id); err != nil {
		return authUsecase.User{}, err
	}
	res := authUsecase.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		IsActive:     user.IsActive,
		IsDeleted:    user.IsDeleted,
		TokenVersion: user.TokenVersion,
		Roles:        user.Roles,
		PasswordHash: user.PasswordHash,
		DateCreated:  user.DateCreated,
		DateUpdated:  user.DateUpdated,
	}
	return res, nil
}

func (s *Storer) Update(ctx context.Context, u authUsecase.User) error {
	ctx, span := web.AddSpan(ctx, "provider.auth.update")
	defer span.End()
	user := StorerUser{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		IsActive:     u.IsActive,
		IsDeleted:    u.IsDeleted,
		TokenVersion: u.TokenVersion,
		Roles:        u.Roles,
		PasswordHash: u.PasswordHash,
		DateCreated:  u.DateCreated,
		DateUpdated:  u.DateUpdated,
	}
	q := `UPDATE users SET
	name = :name, email = :email, is_active = :is_active,
	token_version = :token_version,
	roles = :roles, password_hash = :password_hash,
	date_updated = :date_updated
	WHERE user_id = :user_id;`
	_, err := s.repo.NamedExecContext(ctx, q, user)
	return err
}

func (s *Storer) StoreResetToken(ctx context.Context, rt authUsecase.ResetToken) error {
	ctx, span := web.AddSpan(ctx, "provider.auth.store-reset-token")
	defer span.End()
	token := StorerResetToken{
		TokenID:   rt.TokenID,
		UserID:    rt.UserID,
		Email:     rt.Email,
		IssuedAt:  rt.IssuedAt,
		ExpiresAt: rt.ExpiresAt,
	}
	q := `INSERT INTO reset_tokens (token_id, user_id, email, expires_at, issued_at)
	VALUES (:token_id, :user_id, :email, :expires_at, :issued_at);
	`
	_, err := s.repo.NamedExecContext(ctx, q, token)
	return err
}

func (s *Storer) DeleteResetTokensByUserID(ctx context.Context, uID string) error {
	ctx, span := web.AddSpan(ctx, "provider.auth.delete-reset-tokens-by-user-id")
	defer span.End()
	q := `DELETE from reset_tokens WHERE user_id = $1;`
	_, err := s.repo.ExecContext(ctx, q, uID)
	return err
}

func (s *Storer) QueryResetTokenByID(ctx context.Context, t string) (authUsecase.ResetToken, error) {
	ctx, span := web.AddSpan(ctx, "provider.auth.query-reset-token-by-id")
	defer span.End()
	token := StorerResetToken{}
	q := `SELECT * FROM reset_tokens WHERE token_id = $1`
	if err := s.repo.GetContext(ctx, &token, q, t); err != nil {
		return authUsecase.ResetToken{}, err
	}
	res := authUsecase.ResetToken{
		TokenID:   token.TokenID,
		UserID:    token.UserID,
		Email:     token.Email,
		IssuedAt:  token.IssuedAt,
		ExpiresAt: token.ExpiresAt,
	}
	return res, nil
}

func (s *Storer) DeleteToken(ctx context.Context, t authUsecase.DeleteToken) error {
	ctx, span := web.AddSpan(ctx, "provider.auth.delete-token")
	defer span.End()
	token := StorerDeleteToken{
		TokenID:      t.TokenID,
		Subject:      t.Subject,
		TokenVersion: t.TokenVersion,
		IssuedAt:     t.IssuedAt,
		ExpiresAt:    t.ExpiresAt,
		RevokedAt:    t.RevokedAt,
	}
	q := `INSERT INTO revoked_tokens (token_id, subject, token_version, issued_at, expires_at, revoked_at) 
	VALUES (:token_id, :subject, :token_version, :issued_at, :expires_at, :revoked_at)`
	_, err := s.repo.NamedExecContext(ctx, q, token)
	return err
}
