package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/f4mk/api/internal/app/usecase/user"
	"github.com/f4mk/api/pkg/web"
	"github.com/jmoiron/sqlx"
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
	core *user.Core
	log  *zerolog.Logger
}

func NewService(l *zerolog.Logger, db *sqlx.DB) *UserService {

	repo := NewRepo(db, l)
	core := user.NewCore(repo, l)

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
			getResponseErrorFromUsecase(err),
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
			getResponseErrorFromUsecase(err),
		)
	}
	return web.Respond(ctx, w, res, http.StatusOK)
}

func (us *UserService) CreateUser(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {

	u := user.NewUser{}

	if err := web.Decode(r, &u); err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	res, err := us.core.Create(ctx, u)

	if err != nil {
		return fmt.Errorf(
			"cannot create users: %w",
			getResponseErrorFromUsecase(err),
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

	u := user.UpdateUser{}
	if err := web.Decode(r, &u); err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	res, err := us.core.Update(ctx, id, u)

	if err != nil {
		return fmt.Errorf(
			"cannot update user: %w",
			getResponseErrorFromUsecase(err),
		)
	}
	return web.Respond(ctx, w, res, http.StatusOK)
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
			getResponseErrorFromUsecase(err),
		)
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

func getResponseErrorFromUsecase(err error) error {

	switch {
	case errors.Is(err, user.ErrNotFound):
		return web.NewRequestError(
			err,
			http.StatusNotFound,
		)

	case errors.Is(err, user.ErrForbidden):
		return web.NewRequestError(
			err,
			http.StatusForbidden,
		)

	case errors.Is(err, user.ErrAlreadyExists):
		return web.NewRequestError(
			err,
			http.StatusConflict,
		)
	case errors.Is(err, user.ErrAuthFailed):
		return web.NewRequestError(
			err,
			http.StatusUnauthorized,
		)
	default:
		return err
	}
}

func GetJSONTagName(s interface{}, fieldName string) (string, error) {
	rt := reflect.TypeOf(s)
	field, found := rt.FieldByName(fieldName)
	if !found {
		return "", fmt.Errorf("field %s not found in the struct", fieldName)
	}
	return field.Tag.Get("json"), nil
}
