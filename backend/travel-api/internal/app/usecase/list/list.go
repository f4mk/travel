package list

import (
	"context"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/rs/zerolog"
)

type storer interface {
	// TODO:
	QueryListByID(ctx context.Context, userID string, listID string) (List, error)
	QueryListsByUserID(ctx context.Context, userID string) ([]List, error)
	QueryItemsByListID(ctx context.Context, userID string, listID string) ([]Item, error)
	QueryItemByID(ctx context.Context, userID string, listID string, itemID string) (Item, error)
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

func (c *Core) GetAllLists(ctx context.Context, userID string) ([]List, error) {
	ls, err := c.storer.QueryListsByUserID(ctx, userID)
	if err != nil {
		c.log.Err(err).Msgf("lists: query all: %s", database.ErrQueryDB.Error())
		return nil, database.WrapStorerError(err)
	}
	return ls, nil
}

func (c *Core) GetListByID(ctx context.Context, userID string, listID string) (List, error) {
	l, err := c.storer.QueryListByID(ctx, userID, listID)
	if err != nil {
		c.log.Err(err).Msgf("lists: query: %s", database.ErrQueryDB.Error())
		return List{}, database.WrapStorerError(err)
	}
	return l, nil
}

func (c *Core) GetItemsByListID(ctx context.Context, userID string, listID string) ([]Item, error) {
	is, err := c.storer.QueryItemsByListID(ctx, userID, listID)
	if err != nil {
		c.log.Err(err).Msgf("lists: query: %s", database.ErrQueryDB.Error())
		return []Item{}, database.WrapStorerError(err)
	}
	return is, nil
}

func (c *Core) GetItemByID(ctx context.Context, userID, listID, itemID string) (Item, error) {
	it, err := c.storer.QueryItemByID(ctx, userID, listID, itemID)
	if err != nil {
		c.log.Err(err).Msgf("lists: query: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	return it, nil
}

// TODO: remove when done
//
//revive:disable

func (c *Core) CreateNewList(ctx context.Context, list NewList) (List, error) {
	// Business logic for creating a new list, then:
	return List{}, nil
}

func (c *Core) UpdateListByID(ctx context.Context, ul UpdateList) (List, error) {
	// Business logic for updating a list by ID, then:
	return List{}, nil
}

func (c *Core) DeleteListByID(ctx context.Context, userID string, listID string) error {
	return nil
}

func (c *Core) CreateNewItem(ctx context.Context, userID string, item NewItem) (Item, error) {
	// Business logic for creating a new list, then:
	return Item{}, nil
}

func (c *Core) UpdateItemByID(ctx context.Context, userID string, ui UpdateItem) (Item, error) {
	// Business logic for updating a list by ID, then:
	return Item{}, nil
}

func (c *Core) DeleteItemByID(ctx context.Context, userID string, listID string, itemID string) error {
	return nil
}

// ... Repeat similar methods for items ...
