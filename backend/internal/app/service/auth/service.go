package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	authUsecase "github.com/f4mk/api/internal/app/usecase/auth"
	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/internal/pkg/database"
	"github.com/f4mk/api/pkg/web"

	"github.com/rs/zerolog"
)

type AuthService struct {
	log  *zerolog.Logger
	auth *auth.Auth
	core *authUsecase.Core
}

func NewService(l *zerolog.Logger, auth *auth.Auth, repo authUsecase.Storer) *AuthService {

	core := authUsecase.NewCore(repo, l)

	return &AuthService{
		log:  l,
		auth: auth,
		core: core,
	}
}

func (as *AuthService) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	u := LoginUserDTO{}

	if err := web.Decode(r, &u); err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	au := authUsecase.LoginUser{
		Email:    u.Email,
		Password: u.Password,
	}

	res, err := as.core.Login(ctx, au)

	if err != nil {
		return database.GetResponseErrorFromBusiness(err)
	}

	c := auth.Claims{}
	c.Subject = res.ID
	c.Roles = res.Roles

	newAuthToken, err := as.auth.GenerateToken(c, as.auth.AuthDuration)
	if err != nil {
		return fmt.Errorf("error generating token: %w", err)
	}

	newRefreshToken, err := as.auth.GenerateToken(c, as.auth.RefreshDuration)
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

func (as *AuthService) Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		w.Header().Del("Authorization")
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	claims, err := as.auth.ValidateRefreshToken(ctx, refreshToken.Value)
	if err != nil {

		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	dt := authUsecase.DeleteToken{
		TokenID:   claims.ID,
		Subject:   claims.Subject,
		IssuedAt:  claims.IssuedAt.UTC(),
		ExpiresAt: claims.ExpiresAt.UTC(),
		RevokedAt: time.Now().UTC(),
	}

	if err := as.core.Logout(ctx, dt); err != nil {
		return database.GetResponseErrorFromBusiness(err)
	}

	//try to mark token as revoked
	if err := as.auth.MarkTokenAsRevoked(ctx, auth.TokenParams{
		TokenID:   dt.TokenID,
		Subject:   dt.Subject,
		IssuedAt:  dt.IssuedAt,
		ExpiresAt: dt.ExpiresAt,
		RevokedAt: dt.RevokedAt,
	}); err != nil {
		return err
	}

	//delete all auth info after revoking
	w.Header().Del("Authorization")
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	return web.Respond(ctx, w, nil, http.StatusOK)
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
