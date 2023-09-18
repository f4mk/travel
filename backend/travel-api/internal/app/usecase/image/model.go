package image

import "time"

type Status string

const (
	Pending Status = "pending"
	Loaded  Status = "loaded"
	Deleted Status = "deleted"
)

type Image struct {
	ID          string
	ListID      string
	UserID      string
	ItemID      string
	Private     bool
	Description string
	Status      Status
	DateCreated time.Time
}
