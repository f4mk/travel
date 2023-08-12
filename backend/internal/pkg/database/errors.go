package database

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/f4mk/api/pkg/web"
	"github.com/lib/pq"
)

// lib/pq errorCodeNames
// https://github.com/lib/pq/blob/master/error.go#L178
const (
	uniqueViolation = pq.ErrorCode("23505")
)

var (
	ErrNotFound      = errors.New("not found")
	ErrForbidden     = errors.New("not allowed")
	ErrAuthFailed    = errors.New("authentication failed")
	ErrAlreadyExists = errors.New("already exists")
)

func WrapBusinessError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolation {
		return ErrAlreadyExists
	}

	return err
}

func GetResponseErrorFromBusiness(err error) error {

	switch {
	case errors.Is(err, ErrNotFound):
		return web.NewRequestError(
			err,
			http.StatusNotFound,
		)

	case errors.Is(err, ErrForbidden):
		return web.NewRequestError(
			err,
			http.StatusForbidden,
		)

	case errors.Is(err, ErrAlreadyExists):
		return web.NewRequestError(
			err,
			http.StatusConflict,
		)
	case errors.Is(err, ErrAuthFailed):
		return web.NewRequestError(
			err,
			http.StatusUnauthorized,
		)
	default:
		return err
	}
}
