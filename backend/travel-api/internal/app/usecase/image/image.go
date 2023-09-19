package image

import (
	"context"
	"io"
	"time"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Server interface {
	ServeFile(ctx context.Context, fileID string) ([]byte, error)
	SaveFiles(ctx context.Context, filesID []string, streams []io.Reader) error
	DeleteFiles(ctx context.Context, filesID []string) error
}
type Storer interface {
	QueryByID(ctx context.Context, listID string) (Image, error)
	Create(ctx context.Context, images []Image) error
}
type Converter interface {
	Convert(ctx context.Context, images []io.Reader) ([]io.Reader, error)
}

type Core struct {
	server    Server
	storer    Storer
	converter Converter
	log       *zerolog.Logger
}

func NewCore(l *zerolog.Logger, sr Server, st Storer, cv Converter) *Core {
	return &Core{
		server:    sr,
		storer:    st,
		converter: cv,
		log:       l,
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

func (c *Core) StoreImages(
	ctx context.Context,
	imageStreams []io.Reader,
	itemID, listID, userID string,
) ([]string, error) {
	var imageIDs []string
	var imageItems []Image
	imgStreams, err := c.converter.Convert(ctx, imageStreams)
	if err != nil {
		c.log.Error().Msgf("image: store: convert: %s", err.Error())
		return nil, err
	}
	for range imgStreams {
		img := Image{
			ID:          uuid.New().String(),
			ListID:      listID,
			UserID:      userID,
			ItemID:      itemID,
			Private:     true,
			Description: "",
			Status:      Pending,
			DateCreated: time.Now().UTC(),
		}

		imageItems = append(imageItems, img)
		imageIDs = append(imageIDs, img.ID)
	}

	errCh := make(chan error, 2)
	go func() {
		if err := c.storer.Create(ctx, imageItems); err != nil {
			// no need for cleanup: should be handled by cron
			c.log.Err(err).Msgf("image: create: %s", database.ErrQueryDB.Error())
			errCh <- database.WrapStorerError(err)
		}
		errCh <- nil
	}()
	go func() {
		if err := c.server.SaveFiles(ctx, imageIDs, imgStreams); err != nil {
			c.log.Err(err).Msgf("image: create: save: %s", err.Error())
			// cleanup image storage
			if err := c.server.DeleteFiles(ctx, imageIDs); err != nil {
				c.log.Err(err).Msgf("image: create: rollback: %s", err.Error())
			}
			errCh <- err
		}
		errCh <- nil
	}()
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			return nil, err
		}
	}
	return imageIDs, nil
}
