package image

import (
	"context"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

type Server interface {
	ServeFile(ctx context.Context, fileID string) ([]byte, error)
}
type Storer interface {
	QueryByID(ctx context.Context, listID string) (Image, error)
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

func (c *Core) GetImageByID(ctx context.Context, fileID string, userID string) ([]byte, error) {
	im, err := c.storer.QueryByID(ctx, fileID)
	if err != nil {
		c.log.Err(err).Msgf("image: query by id: %s", database.ErrQueryDB.Error())
		return []byte{}, database.WrapStorerError(err)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Msgf("image: query by id: %s", auth.ErrGetClaims.Error())
		return []byte{}, auth.ErrGetClaims
	}
	if !claims.Authorize(auth.RoleAdmin) && userID != im.UserID && im.Private {
		c.log.Error().Msgf("image: query by id: %s", web.ErrForbidden.Error())
		return []byte{}, web.ErrForbidden
	}

	return c.server.ServeFile(ctx, fileID)
}

// TODO:
func (c *Core) StoreImage(ctx context.Context, itemID, listID, userID string) (string, error) {

	return "", nil
}
