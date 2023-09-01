package list

import "time"

type List struct {
	ID          string    `db:"list_id"`
	UserID      string    `db:"user_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Private     bool      `db:"private"`
	Favorite    bool      `db:"favorite"`
	Completed   bool      `db:"completed"`
	ItemsID     []string  `db:"items"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

type Item struct {
	ID          string    `db:"item_id"`
	ListID      string    `db:"list_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Address     string    `db:"address"`
	PointID     string    `db:"point_id"`
	ImageLinks  []string  `db:"image_links"`
	Links       []string  `db:"links"`
	Visited     bool      `db:"is_visited"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

type Link struct {
	ID   string `db:"link_id"`
	Name string `db:"name"`
	URL  string `db:"url"`
}

type Point struct {
	ID     string  `db:"point_id"`
	ItemID string  `db:"item_id"`
	Lat    float64 `db:"lat"`
	Lng    float64 `db:"lng"`
}
