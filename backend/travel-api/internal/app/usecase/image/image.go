package image

import (
	"context"
	"io"
	"time"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/images"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Server interface {
	ServeFile(ctx context.Context, fileID string) (io.ReadCloser, error)
	SaveFiles(ctx context.Context, filesID []string, streams []io.ReadCloser) error
	TryDeleteFiles(ctx context.Context, filesID []string) error
}
type Storer interface {
	QueryByID(ctx context.Context, fileID string) (Image, error)
	Create(ctx context.Context, images []Image) error
}
type Converter interface {
	Convert(ctx context.Context, images []io.Reader) ([]io.ReadCloser, error)
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

func (c *Core) GetImageByID(ctx context.Context, fileID string, userID string) (io.ReadCloser, error) {
	tID := web.GetTraceID(ctx)
	img, err := c.storer.QueryByID(ctx, fileID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("image: query by id: %s", database.ErrQueryDB.Error())
		return nil, database.WrapStorerError(err)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("image: query by id: %s", auth.ErrGetClaims.Error())
		return nil, auth.ErrGetClaims
	}
	if !claims.Authorize(auth.RoleAdmin) && userID != img.UserID && img.Private {
		c.log.Error().Str("TraceID", tID).Msgf("image: query by id: %s", web.ErrForbidden.Error())
		return nil, web.ErrForbidden
	}

	return c.server.ServeFile(ctx, fileID)
}

func (c *Core) StoreImages(
	ctx context.Context,
	imageStreams []io.Reader,
	listID, userID string,
) ([]string, error) {
	tID := web.GetTraceID(ctx)
	var imageIDs []string
	var imageItems []Image
	imgStreams, err := c.converter.Convert(ctx, imageStreams)
	if err != nil {
		c.log.Error().Str("TraceID", tID).Msgf("image: store: convert: %s", err.Error())
		return nil, err
	}
	for range imgStreams {
		img := Image{
			ID:          uuid.New().String(),
			ListID:      listID,
			UserID:      userID,
			ItemID:      nil,
			Private:     true,
			Description: "",
			Status:      images.Pending,
			DateCreated: time.Now().UTC(),
		}
		imageItems = append(imageItems, img)
		imageIDs = append(imageIDs, img.ID)
	}

	if err := c.storer.Create(ctx, imageItems); err != nil {
		// no need for cleanup: should be handled by cron
		c.log.Err(err).Str("TraceID", tID).Msgf("image: create: %s", database.ErrQueryDB.Error())
		return nil, database.WrapStorerError(err)
	}
	if err := c.server.SaveFiles(ctx, imageIDs, imgStreams); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("image: create: save: %s", err.Error())
		// cleanup image storage
		if err := c.server.TryDeleteFiles(ctx, imageIDs); err != nil {
			c.log.Err(err).Str("TraceID", tID).Msgf("image: create: rollback: %s", err.Error())
			// TODO: store failed ids somewhere
		}
		return nil, err
	}

	return imageIDs, nil
}
