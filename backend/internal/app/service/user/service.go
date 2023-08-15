package user

import (
	"context"
	"fmt"
	"net/http"

	userUsecase "github.com/f4mk/api/internal/app/usecase/user"
	"github.com/f4mk/api/pkg/web"
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
}

func NewService(l *zerolog.Logger, repo userUsecase.Storer) *Service {
	core := userUsecase.NewCore(repo, l)
	return &Service{
		core: core,
		log:  l,
	}
}

func (s *Service) GetUsers(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	res, err := s.core.QueryAll(ctx)
	if err != nil {
		s.log.Err(err).Msg(ErrGetUsersBusiness.Error())
		return fmt.Errorf(
			"cannot get users: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	return web.Respond(ctx, w, res, http.StatusOK)
}

func (s *Service) GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	if err := web.ValidateUUID(id); err != nil {
		s.log.Err(err).Msg(ErrValidateUserID.Error())
		return web.NewRequestError(
			fmt.Errorf("invalid id: %w", err),
			http.StatusBadRequest,
		)
	}
	res, err := s.core.QueryByID(ctx, id)
	if err != nil {
		s.log.Err(err).Msg(ErrGetUserBusiness.Error())
		return fmt.Errorf(
			"cannot get user: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	return web.Respond(ctx, w, res, http.StatusOK)
}

func (s *Service) CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	u := NewUser{}
	if err := web.Decode(r, &u); err != nil {
		s.log.Err(err).Msg(ErrCreateValidate.Error())
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
	res, err := s.core.Create(ctx, nu)
	if err != nil {
		s.log.Err(err).Msg(ErrCreateBusiness.Error())
		return fmt.Errorf(
			"cannot create users: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	return web.Respond(ctx, w, res, http.StatusCreated)
}

func (s *Service) UpdateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	if err := web.ValidateUUID(id); err != nil {
		s.log.Err(err).Msg(ErrUpdateValidateUUID.Error())
		return web.NewRequestError(
			fmt.Errorf("invalid id: %w", err),
			http.StatusBadRequest,
		)
	}
	u := UpdateUser{}
	if err := web.Decode(r, &u); err != nil {
		s.log.Err(err).Msg(ErrUpdateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	uu := userUsecase.UpdateUser{
		Name:  u.Name,
		Email: u.Email,
	}
	res, err := s.core.Update(ctx, id, uu)
	if err != nil {
		s.log.Err(err).Msg(ErrUpdateBusiness.Error())
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
	id := web.Param(r, "id")
	if err := web.ValidateUUID(id); err != nil {
		s.log.Err(err).Msg(ErrDeleteValidateUUID.Error())
		return web.NewRequestError(
			fmt.Errorf("invalid id: %w", err),
			http.StatusBadRequest,
		)
	}
	err := s.core.Delete(ctx, id)
	if err != nil {
		s.log.Err(err).Msg(ErrDeleteBusiness.Error())
		return fmt.Errorf(
			"cannot delete user: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	return web.Respond(ctx, w, nil, http.StatusOK)
}
