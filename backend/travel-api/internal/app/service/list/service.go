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
	ls := []ListResponse{}
	for _, list := range res {
		l := populateListResponse(list)
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
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrListValidateListUUID.Error())
		return err
	}
	res, err := s.core.GetListByID(ctx, claims.Subject, listID)
	if err != nil {
		s.log.Err(err).Msg(ErrGetListsBusiness.Error())
		return fmt.Errorf(
			"cannot query list: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	ls := populateListResponse(res)
	return web.Respond(ctx, w, ls, http.StatusOK)
}

func (s *Service) GetItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrListValidateListUUID.Error())
		return err
	}
	res, err := s.core.GetItemsByListID(ctx, claims.Subject, listID)
	if err != nil {
		s.log.Err(err).Msg(ErrGetListsBusiness.Error())
		return fmt.Errorf(
			"cannot query items: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	is := []ItemResponse{}
	for _, item := range res {
		i := populateItemResponse(item)
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
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrListValidateListUUID.Error())
		return err
	}
	itemID, err := getItemIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrListValidateItemUUID.Error())
		return err
	}
	res, err := s.core.GetItemByID(ctx, claims.Subject, listID, itemID)
	if err != nil {
		s.log.Err(err).Msg(ErrGetListsBusiness.Error())
		return fmt.Errorf(
			"cannot query items: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	i := populateItemResponse(res)
	return web.Respond(ctx, w, i, http.StatusOK)
}

func (s *Service) CreateList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	nl := NewList{}
	if err := web.Decode(r, &nl); err != nil {
		s.log.Err(err).Msg(ErrListCreateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	l := listUsecase.NewList{
		UserID:      claims.Subject,
		Name:        nl.Name,
		Description: nl.Description,
		Private:     nl.Private,
	}
	res, err := s.core.CreateNewList(ctx, l)
	if err != nil {
		s.log.Err(err).Msg(ErrCreateListBusiness.Error())
		return fmt.Errorf(
			"cannot create list: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	pl := populateListResponse(res)
	return web.Respond(ctx, w, pl, http.StatusCreated)
}

func (s *Service) UpdateList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ul := UpdateList{}
	if err := web.Decode(r, &ul); err != nil {
		s.log.Err(err).Msg(ErrListUpdateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrListValidateListUUID.Error())
		return err
	}
	l := listUsecase.UpdateList{
		ID:          listID,
		UserID:      claims.Subject,
		Name:        ul.Name,
		Description: ul.Description,
		Private:     ul.Private,
		Favorite:    ul.Favorite,
		Completed:   ul.Completed,
		ItemsID:     ul.ItemsID,
	}

	res, err := s.core.UpdateListByID(ctx, l)
	if err != nil {
		s.log.Err(err).Msg(ErrUpdateListBusiness.Error())
		return fmt.Errorf(
			"cannot update list: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	pl := populateListResponse(res)
	return web.Respond(ctx, w, pl, http.StatusOK)
}

func (s *Service) DeleteList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	dl := struct{}{}
	if err := web.Decode(r, &dl); err != nil {
		s.log.Err(err).Msg(ErrListDeleteValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrListValidateListUUID.Error())
		return err
	}
	if err := s.core.DeleteListByID(ctx, claims.Subject, listID); err != nil {
		s.log.Err(err).Msg(ErrDeleteListBusiness.Error())
		return fmt.Errorf(
			"cannot delete list: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	return web.Respond(ctx, w, struct{}{}, http.StatusOK)
}

func (s *Service) CreateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ni := NewItem{}
	if err := web.Decode(r, &ni); err != nil {
		s.log.Err(err).Msg(ErrItemCreateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrItemValidateListUUID.Error())
		return err
	}
	np := listUsecase.NewPoint{
		Lat: ni.Point.Lat,
		Lng: ni.Point.Lng,
	}
	var nl *[]listUsecase.NewLink
	if ni.Links != nil {
		tempLinks := []listUsecase.NewLink{}
		for _, link := range *ni.Links {
			l := listUsecase.NewLink{
				Name: link.Name,
				URL:  link.URL,
			}
			tempLinks = append(tempLinks, l)
		}
		nl = &tempLinks
	}
	i := listUsecase.NewItem{
		ListID:      listID,
		Name:        ni.Name,
		Description: ni.Description,
		Address:     ni.Address,
		Point:       np,
		ImageLinks:  ni.ImageLinks,
		Links:       nl,
	}
	res, err := s.core.CreateNewItem(ctx, claims.Subject, i)
	if err != nil {
		s.log.Err(err).Msg(ErrCreateItemBusiness.Error())
		return fmt.Errorf(
			"cannot create item: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	pi := populateItemResponse(res)
	return web.Respond(ctx, w, pi, http.StatusCreated)
}

func (s *Service) UpdateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ui := UpdateItem{}
	if err := web.Decode(r, &ui); err != nil {
		s.log.Err(err).Msg(ErrItemUpdateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrItemValidateListUUID.Error())
		return err
	}
	itemID, err := getItemIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrItemValidateItemUUID.Error())
		return err
	}

	var up *listUsecase.UpdatePoint
	if ui.Point != nil {
		up = &listUsecase.UpdatePoint{
			ID:     ui.Point.ID,
			ItemID: itemID,
			Lat:    ui.Point.Lat,
			Lng:    ui.Point.Lng,
		}
	}

	var ul *[]listUsecase.UpdateLink
	if ui.Links != nil {
		tempLinks := []listUsecase.UpdateLink{}
		for _, link := range *ui.Links {
			l := listUsecase.UpdateLink{
				ID:     link.ID,
				ItemID: itemID,
				Name:   link.Name,
				URL:    link.URL,
			}
			tempLinks = append(tempLinks, l)
		}
		ul = &tempLinks
	}
	i := listUsecase.UpdateItem{
		ID:          itemID,
		ListID:      listID,
		Name:        ui.Name,
		Description: ui.Description,
		Address:     ui.Description,
		Point:       up,
		ImageLinks:  ui.ImageLinks,
		Links:       ul,
		Visited:     ui.Visited,
	}
	res, err := s.core.UpdateItemByID(ctx, claims.Subject, i)
	if err != nil {
		s.log.Err(err).Msg(ErrUpdateItemBusiness.Error())
		return fmt.Errorf(
			"cannot update item: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	pi := populateItemResponse(res)
	return web.Respond(ctx, w, pi, http.StatusOK)
}

func (s *Service) DeleteItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	di := struct{}{}
	if err := web.Decode(r, &di); err != nil {
		s.log.Err(err).Msg(ErrItemDeleteValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrItemValidateListUUID.Error())
		return err
	}
	itemID, err := getItemIDParam(r)
	if err != nil {
		s.log.Err(err).Msg(ErrItemValidateItemUUID.Error())
		return err
	}
	if err := s.core.DeleteItemByID(ctx, claims.Subject, listID, itemID); err != nil {
		s.log.Err(err).Msg(ErrDeleteItemBusiness.Error())
		return fmt.Errorf(
			"cannot delete item: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	return web.Respond(ctx, w, struct{}{}, http.StatusOK)
}

func getListIDParam(r *http.Request) (string, error) {
	listID := web.Param(r, "listID")
	if err := web.ValidateUUID(listID); err != nil {
		return "", web.NewRequestError(
			fmt.Errorf("invalid id: %w", err),
			http.StatusBadRequest,
		)
	}
	return listID, nil
}
func getItemIDParam(r *http.Request) (string, error) {
	itemID := web.Param(r, "itemID")
	if err := web.ValidateUUID(itemID); err != nil {
		return "", web.NewRequestError(
			fmt.Errorf("invalid id: %w", err),
			http.StatusBadRequest,
		)
	}
	return itemID, nil
}

func populateListResponse(res listUsecase.List) ListResponse {
	ls := ListResponse{
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
	return ls
}

func populateItemResponse(res listUsecase.Item) ItemResponse {
	l := []LinkResponse{}
	for _, link := range res.Links {
		l = append(l, LinkResponse{
			ID:     link.ID,
			ItemID: link.ItemID,
			Name:   &link.Name,
			URL:    link.URL,
		})
	}
	i := ItemResponse{
		ID:          res.ID,
		ListID:      res.ListID,
		Name:        res.Name,
		Description: &res.Description,
		Address:     &res.Address,
		Point: PointResponse{
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
	return i
}
