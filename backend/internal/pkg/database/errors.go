package database

import (
	"database/sql"
	"errors"

	"github.com/f4mk/api/pkg/web"
	"github.com/lib/pq"
)

// lib/pq errorCodeNames
// https://github.com/lib/pq/blob/master/error.go#L178
const (
	uniqueViolation = pq.ErrorCode("23505")
)

func WrapBusinessError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return web.ErrNotFound
	}

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolation {
		return web.ErrAlreadyExists
	}

	return err
}
