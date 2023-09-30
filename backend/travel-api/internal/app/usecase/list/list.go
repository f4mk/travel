package list

import (
	"context"
	"time"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type storer interface {
	QueryListByID(ctx context.Context, userID string, listID string) (List, error)
	QueryListsByUserID(ctx context.Context, userID string) ([]List, error)
	QueryItemsByListID(ctx context.Context, userID string, listID string) ([]Item, error)
	QueryItemByID(ctx context.Context, itemID string) (Item, error)
	CreateList(ctx context.Context, list List) error
	UpdateList(ctx context.Context, list List) error
	DeleteListAdmin(ctx context.Context, listID string) error
	DeleteList(ctx context.Context, userID string, listID string) error
	CreateItem(ctx context.Context, item Item) error
	UpdateItem(ctx context.Context, item Item, deleteImages []string) error
	DeleteItemAdmin(ctx context.Context, itemID string) error
	DeleteItem(ctx context.Context, userID string, itemID string) error
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
	ctx, span := web.AddSpan(ctx, "usecase.list.get-all-lists")
	defer span.End()
	tID := web.GetTraceID(ctx)
	ls, err := c.storer.QueryListsByUserID(ctx, userID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("lists: query all: %s", database.ErrQueryDB.Error())
		return nil, database.WrapStorerError(err)
	}
	return ls, nil
}

func (c *Core) GetListByID(ctx context.Context, userID string, listID string) (List, error) {
	ctx, span := web.AddSpan(ctx, "usecase.list.get-list-by-id")
	defer span.End()
	tID := web.GetTraceID(ctx)
	list, err := c.storer.QueryListByID(ctx, userID, listID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("lists: query: %s", database.ErrQueryDB.Error())
		return List{}, database.WrapStorerError(err)
	}
	return list, nil
}

func (c *Core) GetItemsByListID(ctx context.Context, userID string, listID string) ([]Item, error) {
	ctx, span := web.AddSpan(ctx, "usecase.list.get-items-by-list-id")
	defer span.End()
	tID := web.GetTraceID(ctx)
	is, err := c.storer.QueryItemsByListID(ctx, userID, listID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("lists: query: %s", database.ErrQueryDB.Error())
		return []Item{}, database.WrapStorerError(err)
	}
	return is, nil
}

func (c *Core) GetItemByID(ctx context.Context, userID, itemID string) (Item, error) {
	ctx, span := web.AddSpan(ctx, "usecase.list.get-item-by-id")
	defer span.End()
	tID := web.GetTraceID(ctx)
	item, err := c.storer.QueryItemByID(ctx, itemID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("item: query: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("item: query: %s", auth.ErrGetClaims.Error())
		return Item{}, auth.ErrGetClaims
	}
	if !claims.Authorize(auth.RoleAdmin) && (item.UserID != userID && item.Private) {
		c.log.Error().Str("TraceID", tID).Msgf("item: query: %s", web.ErrForbidden.Error())
		return Item{}, web.ErrForbidden
	}
	return item, nil
}

func (c *Core) CreateList(ctx context.Context, nl NewList) (List, error) {
	ctx, span := web.AddSpan(ctx, "usecase.list.create-list")
	defer span.End()
	tID := web.GetTraceID(ctx)
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
		c.log.Err(err).Str("TraceID", tID).Msgf("list: create: %s", database.ErrQueryDB.Error())
		return List{}, database.WrapStorerError(err)
	}
	return list, nil
}

func (c *Core) UpdateList(ctx context.Context, ul UpdateList) (List, error) {
	ctx, span := web.AddSpan(ctx, "usecase.list.update-list")
	defer span.End()
	tID := web.GetTraceID(ctx)
	list, err := c.storer.QueryListByID(ctx, ul.UserID, ul.ID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("list: update: %s", database.ErrQueryDB.Error())
		return List{}, database.WrapStorerError(err)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("list: update: %s", auth.ErrGetClaims.Error())
		return List{}, auth.ErrGetClaims
	}
	if !claims.Authorize(auth.RoleAdmin) && list.UserID != ul.UserID {
		c.log.Error().Str("TraceID", tID).Msgf("item: query: %s", web.ErrForbidden.Error())
		return List{}, web.ErrForbidden
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

	if err := c.storer.UpdateList(ctx, list); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("list: update: %s", database.ErrQueryDB.Error())
		return List{}, database.WrapStorerError(err)
	}
	c.log.Warn().Str("TraceID", tID).Msgf("list: update: %s", list.ID)
	return list, nil
}

func (c *Core) DeleteList(ctx context.Context, userID string, listID string) error {
	ctx, span := web.AddSpan(ctx, "usecase.list.delete-list-by-id")
	defer span.End()
	tID := web.GetTraceID(ctx)
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("list: delete: %s", auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	if claims.Authorize(auth.RoleAdmin) {
		if err := c.storer.DeleteListAdmin(ctx, listID); err != nil {
			c.log.Err(err).Str("TraceID", tID).Msgf("list: delete by admin: %s", database.ErrQueryDB.Error())
			return database.WrapStorerError(err)
		}
		c.log.Warn().Str("TraceID", tID).Msgf("list: deleted by admin: %s", listID)
		return nil
	}
	if err := c.storer.DeleteList(ctx, userID, listID); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("list: delete: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}
	return nil
}

func (c *Core) CreateItem(ctx context.Context, ni NewItem) (Item, error) {
	ctx, span := web.AddSpan(ctx, "usecase.list.create-item")
	defer span.End()
	tID := web.GetTraceID(ctx)
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
		ImagesID:    ni.ImagesID,
		Visited:     false,
		DateCreated: now,
		DateUpdated: now,
	}
	if err := c.storer.CreateItem(ctx, item); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("item: create: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	return item, nil
}

func (c *Core) UpdateItem(ctx context.Context, ui UpdateItem) (Item, error) {
	ctx, span := web.AddSpan(ctx, "usecase.list.update-item")
	defer span.End()
	tID := web.GetTraceID(ctx)
	item, err := c.storer.QueryItemByID(ctx, ui.ID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("item: update: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("item: update: %s", auth.ErrGetClaims.Error())
		return Item{}, auth.ErrGetClaims
	}
	if !claims.Authorize(auth.RoleAdmin) && item.UserID != ui.UserID {
		c.log.Error().Str("TraceID", tID).Msgf("item: update: %s", web.ErrForbidden.Error())
		return Item{}, web.ErrForbidden
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
	var toUpsert []string
	var toDelete []string
	if len(item.ImagesID) != 0 {
		existingMap := make(map[string]bool)
		newMap := make(map[string]bool)

		for _, id := range item.ImagesID {
			existingMap[id] = true
		}
		for _, id := range ui.ImagesID {
			newMap[id] = true
		}
		for id := range existingMap {
			if !newMap[id] {
				toDelete = append(toDelete, id)
			}
		}
		for id := range newMap {
			if !existingMap[id] || existingMap[id] {
				toUpsert = append(toUpsert, id)
			}
		}
	} else {
		toUpsert = ui.ImagesID
	}
	item.ImagesID = toUpsert

	if ui.Visited != nil {
		item.Visited = *ui.Visited
	}
	item.DateUpdated = time.Now().UTC()
	if err := c.storer.UpdateItem(ctx, item, toDelete); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("item: update: %s", database.ErrQueryDB.Error())
		return Item{}, database.WrapStorerError(err)
	}
	c.log.Warn().Str("TraceID", tID).Msgf("item: update: %s", item.ID)
	return item, nil
}

func (c *Core) DeleteItem(ctx context.Context, userID, itemID string) error {
	ctx, span := web.AddSpan(ctx, "usecase.list.delete-item")
	defer span.End()
	tID := web.GetTraceID(ctx)
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("item: delete: %s", auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	if claims.Authorize(auth.RoleAdmin) {
		if err := c.storer.DeleteItemAdmin(ctx, itemID); err != nil {
			c.log.Err(err).Str("TraceID", tID).Msgf("item: delete by admin: %s", database.ErrQueryDB.Error())
			return database.WrapStorerError(err)
		}
		c.log.Warn().Str("TraceID", tID).Msgf("item: deleted by admin: %s", itemID)
		return nil
	}
	if err := c.storer.DeleteItem(ctx, userID, itemID); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("item: delete: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}
	return nil
}
