package list

import (
	"context"
	"time"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type storer interface {
	QueryListByID(ctx context.Context, userID string, listID string) (List, error)
	QueryListsByUserID(ctx context.Context, userID string) ([]List, error)
	QueryItemsByListID(ctx context.Context, userID string, listID string) ([]Item, error)
	QueryItemByID(ctx context.Context, userID string, listID string, itemID string) (Item, error)
	CreateList(ctx context.Context, list List) error
	UpdateListAdmin(ctx context.Context, list List) error
	UpdateList(ctx context.Context, list List) error
	DeleteListAdmin(ctx context.Context, listID string) error
	DeleteList(ctx context.Context, userID string, listID string) error
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

func (c *Core) CreateNewList(ctx context.Context, nl NewList) (List, error) {
	now := time.Now().UTC()
	list := List{
		ID:          uuid.New().String(),
		UserID:      nl.UserID,
		Name:        nl.Name,
		Description: *nl.Description,
		Private:     *nl.Private,
		Completed:   false,
		ItemsID:     []string{},
		DateCreated: now,
		DateUpdated: now,
	}
	if err := c.storer.CreateList(ctx, list); err != nil {
		c.log.Err(err).Msgf("list: create: %s", database.ErrQueryDB.Error())
		return List{}, database.WrapStorerError(err)
	}
	return list, nil
}

func (c *Core) UpdateListByID(ctx context.Context, ul UpdateList) (List, error) {
	l, err := c.storer.QueryListByID(ctx, ul.UserID, ul.ID)
	if err != nil {
		c.log.Err(err).Msgf("list: update: %s", database.ErrQueryDB.Error())
		return List{}, database.WrapStorerError(err)
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Msgf("list: update: %s", auth.ErrGetClaims.Error())
		return List{}, auth.ErrGetClaims
	}
	if ul.Name != nil {
		l.Name = *ul.Name
	}
	if ul.Description != nil {
		l.Description = *ul.Description
	}
	if ul.Private != nil {
		l.Private = *ul.Private
	}
	if ul.Favorite != nil {
		l.Favorite = *ul.Favorite
	}
	if ul.Completed != nil {
		l.Completed = *ul.Completed
	}
	if ul.ItemsID != nil {
		l.ItemsID = *ul.ItemsID
	}
	l.DateUpdated = time.Now().UTC()

	if claims.Authorize(auth.RoleAdmin) {
		if err := c.storer.UpdateListAdmin(ctx, l); err != nil {
			c.log.Err(err).Msgf("list: update: %s", database.ErrQueryDB.Error())
			return List{}, database.WrapStorerError(err)
		}
		c.log.Warn().Msgf("list: update by admin: %s", l.ID)
		return l, nil
	}

	if err := c.storer.UpdateList(ctx, l); err != nil {
		c.log.Err(err).Msgf("list: update: %s", database.ErrQueryDB.Error())
		return List{}, database.WrapStorerError(err)
	}
	return l, nil
}

func (c *Core) DeleteListByID(ctx context.Context, userID string, listID string) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Msgf("list: delete: %s", auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	if claims.Authorize(auth.RoleAdmin) {
		if err := c.storer.DeleteListAdmin(ctx, listID); err != nil {
			c.log.Err(err).Msgf("list: delete by admin: %s", database.ErrQueryDB.Error())
			return database.WrapStorerError(err)
		}
		c.log.Warn().Msgf("list: deleted by admin: %s", listID)
		return nil
	}
	if err := c.storer.DeleteList(ctx, userID, listID); err != nil {
		c.log.Err(err).Msgf("list: delete: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}
	return nil
}

// TODO: remove when done
//
//revive:disable

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
