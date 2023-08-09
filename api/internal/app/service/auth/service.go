package auth

import (
	"context"
	"net/http"

	authUsecase "github.com/f4mk/api/internal/app/usecase/auth"
	"github.com/f4mk/api/internal/pkg/auth"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Service struct {
	log  *zerolog.Logger
	db   *sqlx.DB
	auth *auth.Auth
	core *authUsecase.Core
}

func NewService(l *zerolog.Logger, db *sqlx.DB, auth *auth.Auth) *Service {

	repo := NewRepo(db, l)
	core := authUsecase.NewCore(repo, l)

	return &Service{
		log:  l,
		db:   db,
		auth: auth,
		core: core,
	}
}

// TODO: Implement
func (s *Service) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (s *Service) Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (s *Service) Revoke(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (s *Service) Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (s *Service) PasswordReset(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (s *Service) PasswordChange(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}
