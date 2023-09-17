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
	CreateItem(ctx context.Context, item Item) error
	UpdateItemAdmin(ctx context.Context, item Item) error
	UpdateItem(ctx context.Context, item Item) error
	DeleteItemAdmin(ctx context.Context, itemID string) error
	DeleteItem(ctx context.Context, userID string, listID string, itemID string) error
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
	list, err := c.storer.QueryListByID(ctx, userID, listID)
	if err != nil {
		c.log.Err(err).Msgf("lists: query: %s", database.ErrQueryDB.Error())
		return List{}, database.WrapStorerError(err)
	}
	return list, nil
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
	item, err := c.storer.QueryItemByID(ctx, userID, listID, itemID)
	if err != nil {
		c.log.Err(err).Msgf("lists: query: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	return item, nil
}

func (c *Core) CreateList(ctx context.Context, nl NewList) (List, error) {
	desc := ""
	if nl.Description != nil {
		desc = *nl.Description
	}
	priv := false
	if nl.Private != nil {
		priv = *nl.Private
	}
	now := time.Now().UTC()
	list := List{
		ID:          uuid.New().String(),
		UserID:      nl.UserID,
		Name:        nl.Name,
		Description: desc,
		Private:     priv,
		Completed:   false,
		ItemsID:     nil,
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
	list, err := c.storer.QueryListByID(ctx, ul.UserID, ul.ID)
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
		list.Name = *ul.Name
	}
	if ul.Description != nil {
		list.Description = *ul.Description
	}
	if ul.Private != nil {
		list.Private = *ul.Private
	}
	if ul.Favorite != nil {
		list.Favorite = *ul.Favorite
	}
	if ul.Completed != nil {
		list.Completed = *ul.Completed
	}
	if ul.ItemsID != nil {
		list.ItemsID = ul.ItemsID
	}
	list.DateUpdated = time.Now().UTC()

	if claims.Authorize(auth.RoleAdmin) {
		if err := c.storer.UpdateListAdmin(ctx, list); err != nil {
			c.log.Err(err).Msgf("list: update by admin: %s", database.ErrQueryDB.Error())
			return List{}, database.WrapStorerError(err)
		}
		c.log.Warn().Msgf("list: update by admin: %s", list.ID)
		return list, nil
	}

	if err := c.storer.UpdateList(ctx, list); err != nil {
		c.log.Err(err).Msgf("list: update: %s", database.ErrQueryDB.Error())
		return List{}, database.WrapStorerError(err)
	}
	return list, nil
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

func (c *Core) CreateItem(ctx context.Context, ni NewItem) (Item, error) {
	now := time.Now().UTC()
	itemID := uuid.New().String()
	point := Point{
		ID:     uuid.New().String(),
		ItemID: itemID,
		Lat:    ni.Point.Lat,
		Lng:    ni.Point.Lng,
	}
	item := Item{
		ID:          itemID,
		ListID:      ni.ListID,
		UserID:      ni.UserID,
		Name:        ni.Name,
		Description: ni.Description,
		Address:     ni.Address,
		Point:       point,
		ImageLinks:  ni.ImageLinks,
		Visited:     false,
		DateCreated: now,
		DateUpdated: now,
	}
	if err := c.storer.CreateItem(ctx, item); err != nil {
		c.log.Err(err).Msgf("item: create: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	return item, nil
}

func (c *Core) UpdateItemByID(ctx context.Context, ui UpdateItem) (Item, error) {
	item, err := c.storer.QueryItemByID(ctx, ui.UserID, ui.ListID, ui.ID)
	if err != nil {
		c.log.Err(err).Msgf("item: update: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Msgf("item: update: %s", auth.ErrGetClaims.Error())
		return Item{}, auth.ErrGetClaims
	}
	if ui.Point != nil {
		item.Point.Lat = ui.Point.Lat
		item.Point.Lng = ui.Point.Lng
	}
	if ui.Name != nil {
		item.Name = *ui.Name
	}
	if ui.Description != nil {
		item.Description = ui.Description
	}
	if ui.Address != nil {
		item.Address = ui.Address
	}
	if ui.ImageLinks != nil {
		item.ImageLinks = ui.ImageLinks
	}
	if ui.Visited != nil {
		item.Visited = *ui.Visited
	}
	item.DateUpdated = time.Now().UTC()
	if claims.Authorize(auth.RoleAdmin) {
		if err := c.storer.UpdateItemAdmin(ctx, item); err != nil {
			c.log.Err(err).Msgf("item: update by admin: %s", database.ErrQueryDB.Error())
			return Item{}, database.WrapStorerError(err)
		}
		c.log.Warn().Msgf("item: update by admin: %s", item.ID)
		return item, nil
	}
	// TODO: userID was already checked during QueryItemByID
	// but may be its better to keep it here as well
	if err := c.storer.UpdateItem(ctx, item); err != nil {
		c.log.Err(err).Msgf("item: update: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	return item, nil
}

func (c *Core) DeleteItemByID(ctx context.Context, userID string, listID string, itemID string) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Msgf("item: delete: %s", auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	if claims.Authorize(auth.RoleAdmin) {
		// TODO: may be need to pass listID as well
		if err := c.storer.DeleteItemAdmin(ctx, itemID); err != nil {
			c.log.Err(err).Msgf("item: delete by admin: %s", database.ErrQueryDB.Error())
			return database.WrapStorerError(err)
		}
		c.log.Warn().Msgf("item: deleted by admin: %s", listID)
		return nil
	}
	if err := c.storer.DeleteItem(ctx, userID, listID, itemID); err != nil {
		c.log.Err(err).Msgf("item: delete: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}
	return nil
}
