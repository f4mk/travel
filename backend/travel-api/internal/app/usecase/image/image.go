package image

import (
	"context"

	"github.com/rs/zerolog"
)

type Server interface {
	ServeFile(ctx context.Context, f string) ([]byte, error)
}
type Storer interface {
	QueryListByID(ctx context.Context, listID string) error
}

type Core struct {
	server Server
	storer Storer
	log    *zerolog.Logger
}

func NewCore(l *zerolog.Logger, sr Server, st Storer) *Core {
	return &Core{
		server: sr,
		storer: st,
		log:    l,
	}
}

func (c *Core) QueryByName(ctx context.Context, fname string, userID string) ([]byte, error) {

	return c.server.ServeFile(ctx, fname)

}
