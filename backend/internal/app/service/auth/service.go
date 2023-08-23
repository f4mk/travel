package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	authUsecase "github.com/f4mk/api/internal/app/usecase/auth"
	"github.com/f4mk/api/internal/pkg/auth"
	authPkg "github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/internal/pkg/messages"
	queue "github.com/f4mk/api/pkg/mb"
	"github.com/f4mk/api/pkg/web"

	"github.com/rs/zerolog"
)

type Service struct {
	log  *zerolog.Logger
	auth *authPkg.Auth
	core *authUsecase.Core
	mq   *queue.Channel
}

func NewService(
	l *zerolog.Logger,
	auth *authPkg.Auth,
	storer authUsecase.Storer,
	mq *queue.Channel,
) *Service {

	core := authUsecase.NewCore(storer, l)

	return &Service{
		log:  l,
		auth: auth,
		core: core,
		mq:   mq,
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
	if err := s.auth.StoreUserTokenVersion(ctx, res.UserID, res.TokenVersion); err != nil {
		s.log.Err(err).Msg(ErrLoginStoreTokenVersion.Error())
		return ErrLoginStoreTokenVersion
	}
	c := authPkg.Claims{}
	c.Subject = res.UserID
	c.Roles = res.Roles
	newAuthToken, err := s.auth.GenerateToken(ctx, c, s.auth.AuthDuration)
	if err != nil {
		s.log.Err(err).Msg(ErrLoginGenAuthToken.Error())
		return ErrLoginGenAuthToken
	}
	newRefreshToken, err := s.auth.GenerateToken(ctx, c, s.auth.RefreshDuration)
	if err != nil {
		s.log.Err(err).Msg(ErrLoginGenRefreshToken.Error())
		return ErrLoginGenRefreshToken
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

	c, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	dt := authUsecase.DeleteToken{
		TokenID:      c.ID,
		Subject:      c.Subject,
		TokenVersion: c.TokenVersion,
		IssuedAt:     c.IssuedAt.UTC(),
		ExpiresAt:    c.ExpiresAt.UTC(),
		RevokedAt:    time.Now().UTC(),
	}
	if err := s.core.Logout(ctx, dt); err != nil {
		s.log.Err(err).Msg(ErrLogoutBusiness.Error())
		return web.GetResponseErrorFromBusiness(err)
	}
	if err := s.auth.StoreUserTokenVersion(ctx, c.Subject, c.TokenVersion); err != nil {
		s.log.Err(err).Msg(ErrLoginStoreTokenVersion.Error())
		return ErrLoginStoreTokenVersion
	}
	clearSession(w)
	if err := s.auth.MarkTokenAsRevoked(ctx, authPkg.TokenParams{
		TokenID:      dt.TokenID,
		Subject:      dt.Subject,
		TokenVersion: dt.TokenVersion,
		IssuedAt:     dt.IssuedAt,
		ExpiresAt:    dt.ExpiresAt,
		RevokedAt:    dt.RevokedAt,
	}); err != nil {
		s.log.Err(err).Msg(ErrLogoutRevokeToken.Error())
		return ErrLogoutRevokeToken
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
	c, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	cp := authUsecase.ChangePassword{
		UserID:      c.Subject,
		Password:    p.Password,
		PasswordOld: p.PasswordOld,
	}
	u, err := s.core.ChangePassword(ctx, cp)
	if err != nil {
		s.log.Err(err).Msg(ErrChangePassBusiness.Error())
		return web.GetResponseErrorFromBusiness(err)
	}
	if err := s.auth.StoreUserTokenVersion(ctx, u.UserID, u.TokenVersion); err != nil {
		s.log.Err(err).Msg(ErrLoginStoreTokenVersion.Error())
		return ErrLoginStoreTokenVersion
	}

	return web.Respond(ctx, w, u, http.StatusCreated)
}

func (s *Service) Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	e := struct{}{}
	if err := web.Decode(r, &e); err != nil {
		s.log.Err(err).Msg(ErrLoginDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		s.log.Err(err).Msg(ErrRefreshReadRefreshToken.Error())
		return web.NewRequestError(
			web.ErrAuthFailed,
			http.StatusUnauthorized,
		)
	}
	newClaims, err := s.auth.ValidateToken(ctx, refreshToken.Value)
	if err != nil {
		s.log.Err(err).Msg(ErrRefreshValidateRefreshToken.Error())
		return web.NewRequestError(
			web.ErrAuthFailed,
			http.StatusUnauthorized,
		)
	}
	newAuthToken, err := s.auth.GenerateToken(ctx, newClaims, s.auth.AuthDuration)
	if err != nil {
		s.log.Err(err).Msg(ErrRefreshGenAuthToken.Error())
		return ErrRefreshGenAuthToken
	}
	// trying to revoke old token in case its still valid
	oc, err := s.auth.ParseClaimsFromHeader(r)
	if err != nil {
		s.log.Info().Msg("auth: refresh: previous token not found")
	}
	if oc.ExpiresAt != nil {
		duration := oc.ExpiresAt.Sub(time.Now().UTC())
		if duration > 0 {
			if err := s.auth.MarkTokenAsRevoked(ctx, authPkg.TokenParams{
				TokenID:      oc.ID,
				Subject:      oc.Subject,
				TokenVersion: oc.TokenVersion,
				IssuedAt:     oc.IssuedAt.Time,
				ExpiresAt:    oc.ExpiresAt.Time,
				RevokedAt:    time.Now().UTC(),
			}); err != nil {
				s.log.Err(err).Msg(ErrLogoutRevokeToken.Error())
				return ErrLogoutRevokeToken
			}
		}

	}
	w.Header().Set("Authorization", "Bearer "+newAuthToken)
	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
}

func (s *Service) PasswordReset(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rp := ResetPassword{}
	if err := web.Decode(r, &rp); err != nil {
		s.log.Err(err).Msg(ErrResetPasswordDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	res, err := s.core.ResetPasswordRequest(ctx, rp.Email)
	if err != nil {
		if errors.Is(err, auth.ErrResetTokenReqLimit) {
			return web.NewRequestError(
				err,
				http.StatusTooManyRequests,
			)
		}
		s.log.Err(err).Msg(ErrResetPasswordBusiness.Error())
		//return success to not spoil if the user exists
		if errors.Is(web.ErrNotFound, err) {
			return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
		}
		return fmt.Errorf(
			"cannot reset password: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	m := messages.ResetEmail{
		Email:      res.Email,
		Name:       res.Name,
		ResetToken: res.ResetToken,
	}
	err = s.mq.Publish(ctx, m)
	if err != nil {
		s.log.Err(err).Msg(ErrResetPasswordSendMessage.Error())
		return fmt.Errorf("cannot send message: %w", err)
	}
	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
}

func (s *Service) PasswordResetSubmit(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	srp := SubmitResetPassword{}
	if err := web.Decode(r, &srp); err != nil {
		s.log.Err(err).Msg(ErrValidateResetPasswordDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	sp := authUsecase.SubmitPassword{
		ResetToken: srp.Token,
		Password:   srp.Password,
	}
	u, err := s.core.ResetPasswordSubmit(ctx, sp)
	if err != nil {
		s.log.Err(err).Msg(ErrValidateResetPassword.Error())
		if errors.Is(err, web.ErrNotFound) || errors.Is(err, auth.ErrValidateResetToken) {
			//return forbidden to not spoil if token even exists
			return web.NewRequestError(
				web.ErrForbidden,
				http.StatusForbidden,
			)
		}
		return fmt.Errorf("cannot update password: %w", err)
	}
	if err := s.auth.StoreUserTokenVersion(ctx, u.UserID, u.TokenVersion); err != nil {
		s.log.Err(err).Msg(ErrLoginStoreTokenVersion.Error())
		return ErrLoginStoreTokenVersion
	}
	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
}

func (s *Service) LogoutAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	lu := struct{}{}
	if err := web.Decode(r, &lu); err != nil {
		s.log.Err(err).Msg(ErrLogoutDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	c, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	dt := authUsecase.DeleteToken{
		TokenID:      c.ID,
		Subject:      c.Subject,
		TokenVersion: c.TokenVersion,
		IssuedAt:     c.IssuedAt.UTC(),
		ExpiresAt:    c.ExpiresAt.UTC(),
		RevokedAt:    time.Now().UTC(),
	}
	tv, err := s.core.LogoutAll(ctx, dt)
	if err != nil {
		s.log.Err(err).Msg(ErrLogoutBusiness.Error())
		return web.GetResponseErrorFromBusiness(err)
	}
	if err := s.auth.StoreUserTokenVersion(ctx, c.Subject, tv); err != nil {
		s.log.Err(err).Msg(ErrLoginStoreTokenVersion.Error())
		return ErrLoginStoreTokenVersion
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusOK)
}

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
