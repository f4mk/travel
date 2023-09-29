package user

import (
	"context"
	"fmt"
	"net/http"

	queue "github.com/f4mk/travel/backend/pkg/mb"
	userUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/user"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	authPkg "github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/messages"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

// Service is mediator between upper transfer layer and lower data layer.
// Service directly interacts with the entity and is responsible for
// providing it with the necessary data.
// If that data does not come from a user request, then Service should implement
// an abstraction of a remote data source through its dependency.
// Service injects storer dependency into the data layer.

// Service also aggregates all related handlers under a single name.
// Handlers should always accept a default set of arguments (context, responseWriter, *request),
// send a web response if succeeded, and return an error if failed.
// Errors are supposed to be handled later by the middleware.
// Handlers should handle user input validation as well as convert request data to
// that is acceptable by an underlying usecase

type Service struct {
	core *userUsecase.Core
	log  *zerolog.Logger
	auth *authPkg.Auth
	mq   *queue.Channel
}

func NewService(l *zerolog.Logger, a *authPkg.Auth, c *userUsecase.Core, mq *queue.Channel) *Service {
	return &Service{
		core: c,
		log:  l,
		auth: a,
		mq:   mq,
	}
}

// TODO: should be in admin space or removed
func (s *Service) GetUsers(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	tID := web.GetTraceID(ctx)
	res, err := s.core.QueryAll(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrGetUsersBusiness.Error())
		return fmt.Errorf(
			"cannot get users: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	us := []UserResponse{}
	for _, u := range res {
		us = append(us, UserResponse{
			ID:          u.ID,
			Name:        u.Name,
			Email:       u.Email,
			DateCreated: u.DateCreated,
		})
	}
	return web.Respond(ctx, w, us, http.StatusOK)
}

func (s *Service) GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	id := web.Param(r, "id")
	if err := web.ValidateUUID(id); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrValidateUserID.Error())
		return web.NewRequestError(
			fmt.Errorf("invalid id: %w", err),
			http.StatusBadRequest,
		)
	}
	res, err := s.core.QueryByID(ctx, id)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrGetUserBusiness.Error())
		return fmt.Errorf(
			"cannot get user: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	u := UserResponse{
		ID:          res.ID,
		Name:        res.Name,
		Email:       res.Email,
		DateCreated: res.DateCreated,
	}
	return web.Respond(ctx, w, u, http.StatusOK)
}

func (s *Service) GetMe(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	tID := web.GetTraceID(ctx)
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	res, err := s.core.QueryByID(ctx, claims.Subject)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrGetUserBusiness.Error())
		return fmt.Errorf(
			"cannot get user: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	u := UserResponse{
		ID:          res.ID,
		Name:        res.Name,
		Email:       res.Email,
		DateCreated: res.DateCreated,
	}
	return web.Respond(ctx, w, u, http.StatusOK)
}

func (s *Service) CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	u := NewUser{}
	if err := web.Decode(r, &u); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrCreateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	nu := userUsecase.NewUser{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
	user, token, err := s.core.Create(ctx, nu)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrCreateBusiness.Error())
		return fmt.Errorf(
			"cannot create users: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	m := messages.Message{
		ID:    tID,
		Email: user.Email,
		Name:  user.Name,
		Token: token,
		Type:  messages.RegisterVerify,
	}
	err = s.mq.Publish(ctx, m)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrCreateSendMessage.Error())
		return fmt.Errorf("cannot send message: %w", err)
	}
	cu := UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		DateCreated: user.DateCreated,
	}
	return web.Respond(ctx, w, cu, http.StatusCreated)
}

func (s *Service) VerifyUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	vu := VerifyUser{}
	if err := web.Decode(r, &vu); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrVerifyValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	uu := userUsecase.VerifyUser{
		Email: vu.Email,
		Token: vu.Token,
	}
	res, err := s.core.Verify(ctx, uu)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrVerifyBusiness.Error())
		return fmt.Errorf(
			"cannot verify user: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	ur := UserResponse{
		ID:          res.ID,
		Name:        res.Name,
		Email:       res.Email,
		DateCreated: res.DateCreated,
	}
	return web.Respond(ctx, w, ur, http.StatusCreated)
}

func (s *Service) UpdateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	u := UpdateUser{}
	if err := web.Decode(r, &u); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrUpdateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	uu := userUsecase.UpdateUser{
		ID:       claims.Subject,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
	res, err := s.core.Update(ctx, uu)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrUpdateBusiness.Error())
		return fmt.Errorf(
			"cannot update user: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	ur := UserResponse{
		ID:          res.ID,
		Name:        res.Name,
		Email:       res.Email,
		DateCreated: res.DateCreated,
	}
	return web.Respond(ctx, w, ur, http.StatusOK)
}

func (s *Service) DeleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	u := DeleteUser{}
	if err := web.Decode(r, &u); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrUpdateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	du := userUsecase.DeleteUser{
		ID:       claims.Subject,
		Password: u.Password,
	}
	res, err := s.core.Delete(ctx, du)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrDeleteBusiness.Error())
		return fmt.Errorf(
			"cannot delete user: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	if err := s.auth.StoreUserTokenVersion(ctx, res.ID, res.TokenVersion); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrDeleteStoreTokenVersion.Error())
		return ErrDeleteStoreTokenVersion
	}
	return web.Respond(ctx, w, nil, http.StatusOK)
}
