package user

import (
	"context"
	"fmt"
	"net/http"

	userUsecase "github.com/f4mk/api/internal/app/usecase/user"
	"github.com/f4mk/api/internal/pkg/database"
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
// Handlers should handle user input validation as well as convert request data to DTO
// that is acceptable by an underlying usecase

type UserService struct {
	core *userUsecase.Core
	log  *zerolog.Logger
}

func NewService(l *zerolog.Logger, repo userUsecase.Storer) *UserService {

	core := userUsecase.NewCore(repo, l)

	return &UserService{
		core: core,
		log:  l,
	}
}

func (us *UserService) GetUsers(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {

	res, err := us.core.QueryAll(ctx)

	if err != nil {
		return fmt.Errorf(
			"cannot get users: %w",
			database.GetResponseErrorFromBusiness(err),
		)
	}

	return web.Respond(ctx, w, res, http.StatusOK)
}

func (us *UserService) GetUser(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {

	id := web.Param(r, "id")

	//check if id is valid uuid
	if err := web.ValidateUUID(id); err != nil {
		return web.NewRequestError(
			fmt.Errorf("invalid id: %w", err),
			http.StatusBadRequest,
		)
	}

	res, err := us.core.QueryByID(ctx, id)

	if err != nil {
		return fmt.Errorf(
			"cannot get user: %w",
			database.GetResponseErrorFromBusiness(err),
		)
	}
	return web.Respond(ctx, w, res, http.StatusOK)
}

func (us *UserService) CreateUser(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {

	u := NewUserDTO{}

	if err := web.Decode(r, &u); err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	nu := userUsecase.NewUser{
		Name:            u.Name,
		Email:           u.Email,
		Password:        u.Password,
		PasswordConfirm: u.PasswordConfirm,
	}

	res, err := us.core.Create(ctx, nu)

	if err != nil {
		return fmt.Errorf(
			"cannot create users: %w",
			database.GetResponseErrorFromBusiness(err),
		)
	}

	return web.Respond(ctx, w, res, http.StatusOK)
}

func (us *UserService) UpdateUser(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {

	id := web.Param(r, "id")

	//check if id is valid uuid
	if err := web.ValidateUUID(id); err != nil {
		return web.NewRequestError(
			fmt.Errorf("invalid id: %w", err),
			http.StatusBadRequest,
		)
	}

	u := UpdateUserDTO{}
	if err := web.Decode(r, &u); err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	uu := userUsecase.UpdateUser{
		Name:            u.Name,
		Email:           u.Email,
		Password:        u.Password,
		PasswordConfirm: u.PasswordConfirm,
	}

	res, err := us.core.Update(ctx, id, uu)

	if err != nil {
		return fmt.Errorf(
			"cannot update user: %w",
			database.GetResponseErrorFromBusiness(err),
		)
	}

	ur := UserResponseDTO{
		ID:          res.ID,
		Name:        res.Name,
		Email:       res.Email,
		DateCreated: res.DateCreated,
	}

	return web.Respond(ctx, w, ur, http.StatusOK)
}

func (us *UserService) DeleteUser(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {

	id := web.Param(r, "id")

	//check if id is valid uuid
	if err := web.ValidateUUID(id); err != nil {
		return web.NewRequestError(
			fmt.Errorf("invalid id: %w", err),
			http.StatusBadRequest,
		)
	}

	err := us.core.Delete(ctx, id)

	if err != nil {
		return fmt.Errorf(
			"cannot delete user: %w",
			database.GetResponseErrorFromBusiness(err),
		)
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}
