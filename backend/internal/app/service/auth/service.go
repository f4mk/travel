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

func (s *Service) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	lu := LoginUser{}

	if err := web.Decode(r, &lu); err != nil {
		s.log.Err(err).Msg(ErrLoginDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	au := authUsecase.LoginUser{
		Email:    lu.Email,
		Password: lu.Password,
	}

	res, err := s.core.Login(ctx, au)

	if err != nil {
		s.log.Err(err).Msg(ErrLoginBusiness.Error())
		return web.GetResponseErrorFromBusiness(err)
	}

	c := auth.Claims{}
	c.Subject = res.UserID
	c.Roles = res.Roles

	newAuthToken, err := s.auth.GenerateToken(c, s.auth.AuthDuration)
	if err != nil {
		s.log.Err(err).Msg(ErrLoginGenAuthToken.Error())
		return fmt.Errorf("error generating token: %w", err)
	}

	newRefreshToken, err := s.auth.GenerateToken(c, s.auth.RefreshDuration)
	if err != nil {
		s.log.Err(err).Msg(ErrLoginGenRefreshToken.Error())
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
		ID:          res.UserID,
		Email:       res.Email,
		DateCreated: res.DateCreated,
	}

	return web.Respond(ctx, w, u, http.StatusCreated)
}

func (s *Service) Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	lu := struct{}{}

	if err := web.Decode(r, &lu); err != nil {
		s.log.Err(err).Msg(ErrLogoutDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		s.log.Err(err).Msg(ErrLogoutReadRefreshToken.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	claims, err := s.auth.ValidateRefreshToken(ctx, refreshToken.Value)
	if err != nil {
		s.log.Err(err).Msg(ErrLogoutValidateRefreshToken.Error())
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

	if err := s.core.Logout(ctx, dt); err != nil {
		s.log.Err(err).Msg(ErrLogoutBusiness.Error())
		return web.GetResponseErrorFromBusiness(err)
	}

	clearSession(w)

	if err := s.auth.MarkTokenAsRevoked(ctx, auth.TokenParams{
		TokenID:   dt.TokenID,
		Subject:   dt.Subject,
		IssuedAt:  dt.IssuedAt,
		ExpiresAt: dt.ExpiresAt,
		RevokedAt: dt.RevokedAt,
	}); err != nil {
		s.log.Err(err).Msg(ErrLogoutRevokeToken.Error())
		return fmt.Errorf("error marking token s revoked: %w", err)
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusOK)
}

func (s *Service) ChangePassword(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	p := ChangePassword{}
	if err := web.Decode(r, &p); err != nil {
		s.log.Err(err).Msg(ErrChangePassDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	refreshToken, err := r.Cookie("refresh_token")
	s.log.Err(err).Msg(ErrChangePassReadRefreshToken.Error())

	if err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	claims, err := s.auth.ValidateRefreshToken(ctx, refreshToken.Value)
	if err != nil {
		s.log.Err(err).Msg(ErrChangePassValidateRefreshToken.Error())
		return web.NewRequestError(
			err,
			http.StatusUnauthorized,
		)
	}

	clearSession(w)

	cp := authUsecase.ChangePassword{
		UserID:   claims.Subject,
		Password: p.Password,
	}

	u, err := s.core.ChangePassword(ctx, cp)

	if err != nil {
		s.log.Err(err).Msg(ErrChangePassBusiness.Error())
		return web.GetResponseErrorFromBusiness(err)
	}

	if err := s.auth.MarkTokenAsRevoked(ctx, auth.TokenParams{
		TokenID:   claims.ID,
		Subject:   claims.Subject,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		RevokedAt: time.Now().UTC(),
	}); err != nil {
		s.log.Err(err).Msg(ErrChangePassRevokeToken.Error())
		return fmt.Errorf("error marking token s revoked: %w", err)
	}

	return web.Respond(ctx, w, u, http.StatusCreated)
}

// TODO: Implement
//
//revive:disable
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

//revive:enable

func clearSession(w http.ResponseWriter) {
	w.Header().Del("Authorization")
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})
}
