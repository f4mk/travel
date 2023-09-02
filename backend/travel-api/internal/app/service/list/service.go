package list

import (
	"context"
	"fmt"
	"net/http"

	listUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/list"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

type Service struct {
	core *listUsecase.Core
	log  *zerolog.Logger
}

func NewService(l *zerolog.Logger, c *listUsecase.Core) *Service {

	return &Service{
		core: c,
		log:  l,
	}
}

func (s *Service) GetLists(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	res, err := s.core.GetAllLists(ctx, claims.Subject)
	if err != nil {
		s.log.Err(err).Msg(ErrGetListsBusiness.Error())
		return fmt.Errorf(
			"cannot query lists: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	ls := []List{}
	for _, list := range res {
		l := List{
			ID:          list.ID,
			UserID:      list.UserID,
			Name:        list.Name,
			Description: &list.Description,
			Favorite:    list.Favorite,
			Private:     list.Private,
			Completed:   list.Completed,
			ItemsID:     &list.ItemsID,
			DateCreated: list.DateCreated,
			DateUpdated: &list.DateUpdated,
		}
		ls = append(ls, l)
	}
	return web.Respond(ctx, w, ls, http.StatusOK)
}

func (s *Service) GetList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID := web.Param(r, "listID")
	res, err := s.core.GetListByID(ctx, claims.Subject, listID)
	if err != nil {
		s.log.Err(err).Msg(ErrGetListsBusiness.Error())
		return fmt.Errorf(
			"cannot query list: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	ls := List{
		ID:          res.ID,
		UserID:      res.UserID,
		Name:        res.Name,
		Description: &res.Description,
		Favorite:    res.Favorite,
		Private:     res.Private,
		Completed:   res.Completed,
		ItemsID:     &res.ItemsID,
		DateCreated: res.DateCreated,
		DateUpdated: &res.DateUpdated,
	}
	return web.Respond(ctx, w, ls, http.StatusOK)
}

func (s *Service) GetItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID := web.Param(r, "listID")
	res, err := s.core.GetItemsByListID(ctx, claims.Subject, listID)
	if err != nil {
		s.log.Err(err).Msg(ErrGetListsBusiness.Error())
		return fmt.Errorf(
			"cannot query items: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	is := []Item{}
	for _, item := range res {
		l := []Link{}
		for _, link := range item.Links {
			l = append(l, Link{
				ID:     link.ID,
				ItemID: link.ItemID,
				Name:   &link.Name,
				URL:    link.URL,
			})
		}

		i := Item{
			ID:          item.ID,
			ListID:      item.ListID,
			Name:        item.Name,
			Description: &item.Description,
			Address:     &item.Address,
			Point: Point{
				ID:     item.Point.ID,
				ItemID: item.Point.ItemID,
				Lat:    item.Point.Lat,
				Lng:    item.Point.Lng,
			},
			ImageLinks:  &item.ImageLinks,
			Links:       &l,
			Visited:     item.Visited,
			DateCreated: item.DateCreated,
			DateUpdated: item.DateUpdated,
		}
		is = append(is, i)
	}
	return web.Respond(ctx, w, is, http.StatusOK)
}

func (s *Service) GetItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID := web.Param(r, "listID")
	itemID := web.Param(r, "itemID")
	res, err := s.core.GetItemByID(ctx, claims.Subject, listID, itemID)
	if err != nil {
		s.log.Err(err).Msg(ErrGetListsBusiness.Error())
		return fmt.Errorf(
			"cannot query items: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}

	l := []Link{}
	for _, link := range res.Links {
		l = append(l, Link{
			ID:     link.ID,
			ItemID: link.ItemID,
			Name:   &link.Name,
			URL:    link.URL,
		})
	}
	i := Item{
		ID:          res.ID,
		ListID:      res.ListID,
		Name:        res.Name,
		Description: &res.Description,
		Address:     &res.Address,
		Point: Point{
			ID:     res.Point.ID,
			ItemID: res.Point.ItemID,
			Lat:    res.Point.Lat,
			Lng:    res.Point.Lng,
		},
		ImageLinks:  &res.ImageLinks,
		Links:       &l,
		Visited:     res.Visited,
		DateCreated: res.DateCreated,
		DateUpdated: res.DateUpdated,
	}
	return web.Respond(ctx, w, i, http.StatusOK)
}

// TODO: remove when done
//
// CreateList creates a new list.
//
//revive:disable
func (s *Service) CreateList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// UpdateList updates a specific list by its ID.
func (s *Service) UpdateList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// DeleteList deletes a specific list by its ID.
func (s *Service) DeleteList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// CreateItem creates a new item for a specific list.
func (s *Service) CreateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// UpdateItem updates a specific item by its ID for a specific list.
func (s *Service) UpdateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}

// DeleteItem deletes a specific item by its ID for a specific list.
func (s *Service) DeleteItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Implement your code here...
	return nil
}
