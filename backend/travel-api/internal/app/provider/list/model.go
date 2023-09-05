package list

import (
	"time"

	"github.com/lib/pq"
)

type RepoList struct {
	ID          string          `db:"list_id"`
	UserID      string          `db:"user_id"`
	Name        string          `db:"list_name"`
	Description string          `db:"description"`
	Private     bool            `db:"private"`
	Favorite    bool            `db:"favorite"`
	Completed   bool            `db:"completed"`
	ItemsID     *pq.StringArray `db:"items"`
	DateCreated time.Time       `db:"date_created"`
	DateUpdated time.Time       `db:"date_updated"`
}

type RepoItem struct {
	ID          string          `db:"item_id"`
	ListID      string          `db:"list_id"`
	Name        string          `db:"item_name"`
	Description *string         `db:"description"`
	Address     *string         `db:"address"`
	PointID     string          `db:"point"`
	ImageLinks  *pq.StringArray `db:"image_links"`
	Visited     bool            `db:"is_visited"`
	DateCreated time.Time       `db:"date_created"`
	DateUpdated time.Time       `db:"date_updated"`
}

type RepoPoint struct {
	ID     string  `db:"point_id"`
	ItemID string  `db:"item_id"`
	Lat    float64 `db:"lat"`
	Lng    float64 `db:"lng"`
}
