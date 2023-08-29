package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/f4mk/travel/backend/pkg/web"
	"github.com/lib/pq"
)

// lib/pq errorCodeNames
// https://github.com/lib/pq/blob/master/error.go#L178
const (
	uniqueViolation  = pq.ErrorCode("23505")
	deadlineExceeded = pq.ErrorCode("57014")
)

var ErrQueryDB = errors.New("error querying db")

func WrapStorerError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return web.ErrNotFound
	}

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolation {
		return web.ErrAlreadyExists
	}

	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == deadlineExceeded {
			return context.DeadlineExceeded
		}
	}

	return err
}
