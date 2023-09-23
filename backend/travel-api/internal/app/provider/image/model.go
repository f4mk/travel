package image

import (
	"time"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/images"
)

type StorerImage struct {
	ID          string        `db:"image_id"`
	ListID      string        `db:"list_id"`
	UserID      string        `db:"user_id"`
	ItemID      *string       `db:"item_id"`
	Private     bool          `db:"private"`
	Description string        `db:"description"`
	Status      images.Status `db:"status"`
	DateCreated time.Time     `db:"date_created"`
}
