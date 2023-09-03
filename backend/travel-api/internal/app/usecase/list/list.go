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
	CreateItem(ctx context.Context, userID string, item Item) error
	UpdateItemAdmin(ctx context.Context, item Item) error
	UpdateItem(ctx context.Context, userID string, item Item) error
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
		l.ItemsID = ul.ItemsID
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

func (c *Core) CreateItem(ctx context.Context, userID string, ni NewItem) (Item, error) {
	now := time.Now().UTC()
	itemID := uuid.New().String()
	point := Point{
		ID:     uuid.New().String(),
		ItemID: itemID,
		Lat:    ni.Point.Lat,
		Lng:    ni.Point.Lng,
	}
	links := []Link{}
	if ni.Links != nil {
		for _, link := range *ni.Links {
			l := Link{
				ID:     uuid.New().String(),
				ItemID: itemID,
				Name:   link.Name,
				URL:    link.URL,
			}
			links = append(links, l)
		}
	}

	item := Item{
		ID:          itemID,
		ListID:      ni.ListID,
		Name:        ni.Name,
		Description: ni.Description,
		Address:     ni.Address,
		Point:       point,
		ImageLinks:  ni.ImageLinks,
		// TODO: handle all nil dereferencing in one place
		Links:       &links,
		Visited:     false,
		DateCreated: now,
		DateUpdated: now,
	}
	if err := c.storer.CreateItem(ctx, userID, item); err != nil {
		c.log.Err(err).Msgf("item: create: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	return item, nil
}

func (c *Core) UpdateItemByID(ctx context.Context, userID string, ui UpdateItem) (Item, error) {
	i, err := c.storer.QueryItemByID(ctx, userID, ui.ListID, ui.ID)
	if err != nil {
		c.log.Err(err).Msgf("item: update: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Msgf("item: update: %s", auth.ErrGetClaims.Error())
		return Item{}, auth.ErrGetClaims
	}

	if ui.Name != nil {
		i.Name = *ui.Name
	}
	if ui.Description != nil {
		i.Description = ui.Description
	}
	if ui.Address != nil {
		i.Address = ui.Address
	}
	if ui.Point != nil {
		i.Point = Point{
			ID:     ui.Point.ID,
			ItemID: ui.Point.ItemID,
			Lat:    ui.Point.Lat,
			Lng:    ui.Point.Lng,
		}
	}
	if ui.ImageLinks != nil {
		i.ImageLinks = ui.ImageLinks
	}
	if ui.Links != nil {
		links := []Link{}
		for _, link := range *ui.Links {

			l := Link{
				ID:     link.ID,
				ItemID: link.ItemID,
			}
			if link.Name != nil {
				l.Name = link.Name
			}
			if link.URL != nil {
				l.URL = *link.URL
			}
			links = append(links, l)
		}
		i.Links = &links
	}
	if ui.Visited != nil {
		i.Visited = *ui.Visited
	}
	i.DateUpdated = time.Now().UTC()
	if claims.Authorize(auth.RoleAdmin) {
		if err := c.storer.UpdateItemAdmin(ctx, i); err != nil {
			c.log.Err(err).Msgf("item: update: %s", database.ErrQueryDB.Error())
			return Item{}, database.WrapStorerError(err)
		}
		c.log.Warn().Msgf("item: update by admin: %s", i.ID)
		return i, nil
	}
	// TODO: userID was already checked during QueryItemByID
	// but may be its better to keep it here as well
	if err := c.storer.UpdateItem(ctx, userID, i); err != nil {
		c.log.Err(err).Msgf("item: update: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	return i, nil
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
