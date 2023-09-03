package list

import "time"

type List struct {
	ID          string
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

type NewList struct {
	UserID      string
	Name        string
	Description *string
	Private     *bool
}

type UpdateList struct {
	ID          string
	UserID      string
	Name        *string
	Description *string
	Private     *bool
	Favorite    *bool
	Completed   *bool
	ItemsID     *[]string
}

type Item struct {
	ID          string
	ListID      string
	Name        string
	Description string
	Address     string
	Point       Point
	ImageLinks  []string
	Links       []Link
	Visited     bool
	DateCreated time.Time
	DateUpdated time.Time
}

type Link struct {
	ID     string
	ItemID string
	Name   string
	URL    string
}

type Point struct {
	ID     string
	ItemID string
	Lat    float64
	Lng    float64
}

type NewItem struct {
	ListID      string
	Name        string
	Description *string
	Address     *string
	Point       NewPoint
	ImageLinks  []string
	Links       []NewLink
}

type NewLink struct {
	Name *string
	URL  string
}

type NewPoint struct {
	Lat float64
	Lng float64
}

type UpdateItem struct {
	ID          string
	ListID      string
	Name        *string
	Description *string
	Address     *string
	Point       *UpdatePoint
	ImageLinks  *[]string
	Links       *[]UpdateLink
	Visited     *bool
}

type UpdateLink struct {
	ID     string
	ItemID string
	Name   *string
	URL    *string
}

type UpdatePoint struct {
	ID     string
	ItemID string
	Lat    float64
	Lng    float64
}
