package image

import (
	"context"

	"github.com/f4mk/travel/backend/travel-api/internal/app/usecase/image"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
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

func (s *Storer) QueryByID(ctx context.Context, imageID string) (image.Image, error) {
	ctx, span := web.AddSpan(ctx, "provider.image.storer.query-by-id")
	defer span.End()
	img := StorerImage{}
	q := `SELECT * from images where imageID=$1`
	if err := s.repo.SelectContext(ctx, &img, q, imageID); err != nil {
		return image.Image{}, err
	}

	res := image.Image{
		ID:          img.ID,
		ListID:      img.ListID,
		UserID:      img.UserID,
		ItemID:      img.ItemID,
		Private:     img.Private,
		Description: img.Description,
		Status:      img.Status,
		DateCreated: img.DateCreated,
	}
	return res, nil
}

func (s *Storer) Create(ctx context.Context, images []image.Image) (err error) {
	ctx, span := web.AddSpan(ctx, "provider.image.storer.create")
	defer span.End()
	q := `INSERT INTO images (
		image_id, list_id, user_id,
		item_id, private,	description,
		status,	date_created,
	) VALUES (
		:image_id, :list_id, :user_id,
		:item_id, :private,	:description,
		:status, :date_created,
	);
	`
	tx, err := s.repo.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				s.log.Err(rErr).Msg("failed to rollback after error")
			}
		}
	}()

	for _, img := range images {
		_, err = tx.NamedExecContext(ctx, q, img)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}
