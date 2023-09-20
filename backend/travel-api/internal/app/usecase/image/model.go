package image

import (
	"time"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/images"
)

type Image struct {
	ID          string
	ListID      string
	UserID      string
	ItemID      *string
	Private     bool
	Description string
	Status      images.Status
	DateCreated time.Time
}
