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

// TODO: implement
func (r *Repo) Update(ctx context.Context, u authUsecase.User) error {
	return nil
}

// TODO: implement
func (r *Repo) Delete(ctx context.Context, t string) error {
	return nil
}

// TODO: implement
func (r *Repo) DeleteAll(ctx context.Context, uID string) error {
	return nil
}

// TODO: implement
func (r *Repo) QueryByEmail(ctx context.Context, uID string) (authUsecase.User, error) {
	return authUsecase.User{}, nil
}
