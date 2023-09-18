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
	ItemsID     *[]string
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
	UserID      string
	Name        string
	Description *string
	Address     *string
	Point       Point
	ImagesID    *[]string
	Visited     bool
	DateCreated time.Time
	DateUpdated time.Time
}

type Point struct {
	ID     string
	ItemID string
	Lat    float64
	Lng    float64
}

type NewItem struct {
	ListID      string
	UserID      string
	Name        string
	Description *string
	Address     *string
	Point       NewPoint
	ImagesID    *[]string
}

type NewPoint struct {
	Lat float64
	Lng float64
}

type UpdateItem struct {
	ID          string
	ListID      string
	UserID      string
	Name        *string
	Description *string
	Address     *string
	Point       *UpdatePoint
	ImagesID    *[]string
	Visited     *bool
}

type UpdatePoint struct {
	Lat float64
	Lng float64
}
