package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	queue "github.com/f4mk/travel/backend/pkg/mb"
	authUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	authPkg "github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/messages"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"

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
	a *authPkg.Auth,
	c *authUsecase.Core,
	mq *queue.Channel,
) *Service {

	return &Service{
		log:  l,
		auth: a,
		core: c,
		mq:   mq,
	}
}

func (s *Service) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	ctx, span := web.AddSpan(ctx, "service.auth.login")
	defer span.End()
	lu := LoginUser{}
	if err := web.Decode(r, &lu); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	au := authUsecase.LoginUser{
		Email:    strings.ToLower(lu.Email),
		Password: lu.Password,
	}
	res, err := s.core.Login(ctx, au)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginBusiness.Error())
		return web.GetResponseErrorFromBusiness(err)
	}
	if err := s.auth.StoreUserTokenVersion(ctx, res.UserID, res.TokenVersion); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginStoreTokenVersion.Error())
		return ErrLoginStoreTokenVersion
	}
	c := authPkg.Claims{}
	c.Subject = res.UserID
	c.Roles = res.Roles
	newAuthToken, err := s.auth.GenerateToken(ctx, c, s.auth.AuthDuration)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginGenAuthToken.Error())
		return ErrLoginGenAuthToken
	}
	newRefreshToken, err := s.auth.GenerateToken(ctx, c, s.auth.RefreshDuration)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginGenRefreshToken.Error())
		return ErrLoginGenRefreshToken
	}
	w.Header().Set("Authorization", "Bearer "+newAuthToken)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
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
	ctx, span := web.AddSpan(ctx, "service.auth.logout")
	defer span.End()
	tID := web.GetTraceID(ctx)
	lu := struct{}{}
	if err := web.Decode(r, &lu); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLogoutDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	c, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
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
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLogoutBusiness.Error())
		return web.GetResponseErrorFromBusiness(err)
	}
	if err := s.auth.StoreUserTokenVersion(ctx, c.Subject, c.TokenVersion); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginStoreTokenVersion.Error())
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
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLogoutRevokeToken.Error())
		return ErrLogoutRevokeToken
	}
	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
}

func (s *Service) ChangePassword(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := web.AddSpan(ctx, "service.auth.change-password")
	defer span.End()
	tID := web.GetTraceID(ctx)
	p := ChangePassword{}
	if err := web.Decode(r, &p); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrChangePassDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	c, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	cp := authUsecase.ChangePassword{
		UserID:      c.Subject,
		Password:    p.Password,
		PasswordOld: p.PasswordOld,
	}
	res, err := s.core.ChangePassword(ctx, cp)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrChangePassBusiness.Error())
		return web.GetResponseErrorFromBusiness(err)
	}
	if err := s.auth.StoreUserTokenVersion(ctx, res.ID, res.TokenVersion); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginStoreTokenVersion.Error())
		return ErrLoginStoreTokenVersion
	}
	u := UserResponse{
		ID:          res.ID,
		Name:        res.Name,
		Email:       res.Email,
		DateCreated: res.DateCreated,
	}

	return web.Respond(ctx, w, u, http.StatusCreated)
}

func (s *Service) Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := web.AddSpan(ctx, "service.auth.refresh")
	defer span.End()
	tID := web.GetTraceID(ctx)
	e := struct{}{}
	if err := web.Decode(r, &e); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrRefreshReadRefreshToken.Error())
		return web.NewRequestError(
			web.ErrAuthFailed,
			http.StatusUnauthorized,
		)
	}
	newClaims, err := s.auth.ValidateToken(ctx, refreshToken.Value)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrRefreshValidateRefreshToken.Error())
		return web.NewRequestError(
			web.ErrAuthFailed,
			http.StatusUnauthorized,
		)
	}
	newAuthToken, err := s.auth.GenerateToken(ctx, newClaims, s.auth.AuthDuration)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrRefreshGenAuthToken.Error())
		return ErrRefreshGenAuthToken
	}
	// trying to revoke old token in case its still valid
	oc, err := s.auth.ParseClaimsFromHeader(r)
	if err != nil {
		s.log.Info().Str("TraceID", tID).Msg("auth: refresh: previous token not found")
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
				s.log.Err(err).Str("TraceID", tID).Msg(ErrLogoutRevokeToken.Error())
				return ErrLogoutRevokeToken
			}
		}

	}
	w.Header().Set("Authorization", "Bearer "+newAuthToken)
	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
}

func (s *Service) Validate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := web.AddSpan(ctx, "service.auth.validate")
	defer span.End()
	tID := web.GetTraceID(ctx)
	e := struct{}{}
	if err := web.Decode(r, &e); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrRefreshReadRefreshToken.Error())
		return web.NewRequestError(
			web.ErrAuthFailed,
			http.StatusUnauthorized,
		)
	}
	_, err = s.auth.ValidateToken(ctx, refreshToken.Value)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrRefreshValidateRefreshToken.Error())
		return web.NewRequestError(
			web.ErrAuthFailed,
			http.StatusUnauthorized,
		)
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
}

func (s *Service) PasswordReset(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := web.AddSpan(ctx, "service.auth.password-reset")
	defer span.End()
	tID := web.GetTraceID(ctx)
	rp := ResetPassword{}
	if err := web.Decode(r, &rp); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrResetPasswordDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	res, err := s.core.ResetPasswordRequest(ctx, strings.ToLower(rp.Email))
	if err != nil {
		if errors.Is(err, auth.ErrResetTokenReqLimit) {
			return web.NewRequestError(
				err,
				http.StatusTooManyRequests,
			)
		}
		s.log.Err(err).Str("TraceID", tID).Msg(ErrResetPasswordBusiness.Error())
		//return success to not spoil if the user exists
		if errors.Is(web.ErrNotFound, err) {
			return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
		}
		return fmt.Errorf(
			"cannot reset password: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	m := messages.Message{
		ID:    tID,
		Email: res.Email,
		Name:  res.Name,
		Token: res.ResetToken,
		Type:  messages.ResetPassword,
	}
	err = s.mq.Publish(ctx, m)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrResetPasswordSendMessage.Error())
		return fmt.Errorf("cannot send message: %w", err)
	}
	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
}

func (s *Service) PasswordResetSubmit(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := web.AddSpan(ctx, "service.auth.password-reset-submit")
	defer span.End()
	tID := web.GetTraceID(ctx)
	srp := SubmitResetPassword{}
	if err := web.Decode(r, &srp); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrValidateResetPasswordDecode.Error())
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
		s.log.Err(err).Str("TraceID", tID).Msg(ErrValidateResetPassword.Error())
		if errors.Is(err, web.ErrNotFound) || errors.Is(err, auth.ErrValidateResetToken) {
			//return forbidden to not spoil if token even exists
			return web.NewRequestError(
				web.ErrForbidden,
				http.StatusForbidden,
			)
		}
		return fmt.Errorf("cannot update password: %w", err)
	}
	if err := s.auth.StoreUserTokenVersion(ctx, u.ID, u.TokenVersion); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginStoreTokenVersion.Error())
		return ErrLoginStoreTokenVersion
	}
	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
}

func (s *Service) LogoutAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := web.AddSpan(ctx, "service.auth.logout-all")
	defer span.End()
	tID := web.GetTraceID(ctx)
	lu := struct{}{}
	if err := web.Decode(r, &lu); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLogoutDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	c, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
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
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLogoutBusiness.Error())
		return web.GetResponseErrorFromBusiness(err)
	}
	if err := s.auth.StoreUserTokenVersion(ctx, c.Subject, tv); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrLoginStoreTokenVersion.Error())
		return ErrLoginStoreTokenVersion
	}
	return web.Respond(ctx, w, struct{}{}, http.StatusCreated)
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
