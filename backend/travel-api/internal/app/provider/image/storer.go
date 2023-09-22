package image

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Storer struct {
	repo *sqlx.DB
	log  *zerolog.Logger
}

func NewStorer(l *zerolog.Logger, r *sqlx.DB) *Storer {
	return &Storer{repo: r, log: l}
}
