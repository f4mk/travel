package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/jackc/pgx/v5/pgconn"
)

// errorCodeNames
// https://www.postgresql.org/docs/current/errcodes-appendix.html

const (
	uniqueViolation  = "23505"
	deadlineExceeded = "57014"
)

var ErrQueryDB = errors.New("error querying db")

func WrapStorerError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return web.ErrNotFound
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case uniqueViolation:
			return web.ErrAlreadyExists
		case deadlineExceeded:
			return context.DeadlineExceeded
		}
	}

	return err
}
