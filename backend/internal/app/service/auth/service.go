package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	authUsecase "github.com/f4mk/api/internal/app/usecase/auth"
	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/pkg/web"

	"github.com/rs/zerolog"
)

type Service struct {
	log  *zerolog.Logger
	auth *auth.Auth
	core *authUsecase.Core
}

func NewService(l *zerolog.Logger, auth *auth.Auth, repo authUsecase.Storer) *Service {

	core := authUsecase.NewCore(repo, l)

	return &Service{
		log:  l,
		auth: auth,
		core: core,
	}
}

func (as *Service) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	lu := LoginUser{}

	if err := web.Decode(r, &lu); err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	au := authUsecase.LoginUser{
		Email:    lu.Email,
		Password: lu.Password,
	}

	res, err := as.core.Login(ctx, au)

	if err != nil {
		return web.GetResponseErrorFromBusiness(err)
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

	u := UserResponse{
		Name:        res.Name,
		ID:          res.ID,
		Email:       res.Email,
		DateCreated: res.DateCreated,
	}

	return web.Respond(ctx, w, u, http.StatusCreated)
}

func (as *Service) Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// TODO: do i need this?
	lu := struct{}{}

	if err := web.Decode(r, &lu); err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

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
			http.StatusUnauthorized,
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
		return web.GetResponseErrorFromBusiness(err)
	}

	//try to mark token as revoked
	if err := as.auth.MarkTokenAsRevoked(ctx, auth.TokenParams{
		TokenID:   dt.TokenID,
		Subject:   dt.Subject,
		IssuedAt:  dt.IssuedAt,
		ExpiresAt: dt.ExpiresAt,
		RevokedAt: dt.RevokedAt,
	}); err != nil {
		return fmt.Errorf("error marking token as revoked: %w", err)
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

	return web.Respond(ctx, w, struct{}{}, http.StatusOK)
}

// TODO: Implement
//
//revive:disable
func (as *Service) Revoke(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (as *Service) Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (as *Service) PasswordReset(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: Implement
func (as *Service) PasswordChange(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

//revive:enable
