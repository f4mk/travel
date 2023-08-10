package auth

import (
	"context"
	"fmt"
	"net/http"

	authUsecase "github.com/f4mk/api/internal/app/usecase/auth"
	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/internal/pkg/database"
	"github.com/f4mk/api/pkg/web"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type AuthService struct {
	log  *zerolog.Logger
	db   *sqlx.DB
	auth *auth.Auth
	core *authUsecase.Core
}

func NewService(l *zerolog.Logger, db *sqlx.DB, auth *auth.Auth) *AuthService {

	repo := NewRepo(db, l)
	core := authUsecase.NewCore(repo, l)

	return &AuthService{
		log:  l,
		db:   db,
		auth: auth,
		core: core,
	}
}

// TODO: Implement
func (as *AuthService) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	u := authUsecase.LoginUser{}

	if err := web.Decode(r, &u); err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	au, err := as.core.Login(ctx, u)

	if err != nil {
		return database.GetResponseErrorFromBusiness(err)
	}

	//so far get the auth user info,
	//need to create token with user subject, set token to user credentials

	c := auth.Claims{}
	c.Subject = au.ID
	c.Roles = au.Roles
	newAuthToken, newRefreshToken, err := as.auth.GenerateTokens(c)

	if err != nil {
		return fmt.Errorf("error generating token: %w", err)
	}

	w.Header().Set("Authorization", "Bearer "+newAuthToken)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	return web.Respond(ctx, w, au, http.StatusOK)
}

// TODO: Implement
func (as *AuthService) Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (as *AuthService) Revoke(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (as *AuthService) Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (as *AuthService) PasswordReset(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (as *AuthService) PasswordChange(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}
