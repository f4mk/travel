package list

import (
	"context"

	"github.com/rs/zerolog"
)

type storer interface {
	// TODO:
	QueryListByID(ctx context.Context, listID string) (ListWithIDs, error)
}

type Core struct {
	storer storer
	log    *zerolog.Logger
}

func NewCore(l *zerolog.Logger, s storer) *Core {
	return &Core{
		storer: s,
		log:    l,
	}
}

// TODO: remove when done
//
//revive:disable
func (c *Core) GetAllLists(ctx context.Context, userID string) error {
	return nil
}

func (c *Core) CreateNewList(ctx context.Context, list NewList) error {
	// Business logic for creating a new list, then:
	return nil
}

func (c *Core) GetListByID(ctx context.Context, listID string) error {
	return nil
}

func (c *Core) UpdateListByID(ctx context.Context, listID string, updatedList UpdateList) error {
	// Business logic for updating a list by ID, then:
	return nil
}

func (c *Core) DeleteListByID(ctx context.Context, listID string) error {
	return nil
}

// ... Repeat similar methods for items ...
