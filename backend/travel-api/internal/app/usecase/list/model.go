package list

import "time"

type ListWithIDs struct {
	ID          string
	Type        string
	UserID      string
	Name        string
	Description string
	Private     bool
	Favorite    bool
	Completed   bool
	ItemsID     []string
	DateCreated time.Time
	DateUpdated time.Time
}

type ListWithItems struct {
	ID          string
	Type        string
	UserID      string
	Name        string
	Description string
	Private     bool
	Favorite    bool
	Completed   bool
	Items       []Item
	DateCreated time.Time
	DateUpdated time.Time
}

type NewList struct {
	UserID      string
	Name        string
	Description *string
	Private     bool
}

type UpdateList struct {
	ID          string
	UserID      string
	Name        *string
	Description *string
	Private     bool
	Favorite    bool
	Completed   bool
	ItemsID     []string
}

type Item struct {
	ID          string
	ListID      string
	Name        string
	Description string
	Address     string
	Point       Point
	ImageLinks  []string
	LinksID     []Link
	Visited     bool
	DateCreated time.Time
	DateUpdated time.Time
}

type Link struct {
	ID   string
	Name string
	URL  string
}

type Point struct {
	ID     string
	ItemID string
	Lat    float64
	Lng    float64
}
