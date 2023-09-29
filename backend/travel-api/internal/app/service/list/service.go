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
	tID := web.GetTraceID(ctx)
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	res, err := s.core.GetAllLists(ctx, claims.Subject)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrGetListsBusiness.Error())
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
	tID := web.GetTraceID(ctx)
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrListValidateListUUID.Error())
		return err
	}
	res, err := s.core.GetListByID(ctx, claims.Subject, listID)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrGetListsBusiness.Error())
		return fmt.Errorf(
			"cannot query list: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	ls := populateListResponse(res)
	return web.Respond(ctx, w, ls, http.StatusOK)
}

func (s *Service) GetItems(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrListValidateListUUID.Error())
		return err
	}
	res, err := s.core.GetItemsByListID(ctx, claims.Subject, listID)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrGetListsBusiness.Error())
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
	tID := web.GetTraceID(ctx)
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	itemID, err := getItemIDParam(r)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrListValidateItemUUID.Error())
		return err
	}
	res, err := s.core.GetItemByID(ctx, claims.Subject, itemID)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrGetListsBusiness.Error())
		return fmt.Errorf(
			"cannot query items: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	i := populateItemResponse(res)
	return web.Respond(ctx, w, i, http.StatusOK)
}

func (s *Service) CreateList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	nl := NewList{}
	if err := web.Decode(r, &nl); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrListCreateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	l := listUsecase.NewList{
		UserID:      claims.Subject,
		Name:        nl.Name,
		Description: nl.Description,
		Private:     nl.Private,
	}
	res, err := s.core.CreateList(ctx, l)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrCreateListBusiness.Error())
		return fmt.Errorf(
			"cannot create list: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	pl := populateListResponse(res)
	return web.Respond(ctx, w, pl, http.StatusCreated)
}

func (s *Service) UpdateList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	ul := UpdateList{}
	if err := web.Decode(r, &ul); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrListUpdateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrListValidateListUUID.Error())
		return err
	}
	var isID []string
	if ul.ItemsID != nil {
		isID = *ul.ItemsID
	}
	l := listUsecase.UpdateList{
		ID:          listID,
		UserID:      claims.Subject,
		Name:        ul.Name,
		Description: ul.Description,
		Private:     ul.Private,
		Favorite:    ul.Favorite,
		Completed:   ul.Completed,
		ItemsID:     isID,
	}

	res, err := s.core.UpdateListByID(ctx, l)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrUpdateListBusiness.Error())
		return fmt.Errorf(
			"cannot update list: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	pl := populateListResponse(res)
	return web.Respond(ctx, w, pl, http.StatusOK)
}

func (s *Service) DeleteList(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	dl := struct{}{}
	if err := web.Decode(r, &dl); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrListDeleteValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrListValidateListUUID.Error())
		return err
	}
	if err := s.core.DeleteListByID(ctx, claims.Subject, listID); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrDeleteListBusiness.Error())
		return fmt.Errorf(
			"cannot delete list: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	return web.Respond(ctx, w, struct{}{}, http.StatusOK)
}

func (s *Service) CreateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	ni := NewItem{}
	if err := web.Decode(r, &ni); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrItemCreateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrItemValidateListUUID.Error())
		return err
	}
	np := listUsecase.NewPoint{
		Lat: ni.Point.Lat,
		Lng: ni.Point.Lng,
	}
	var imgsID []string
	if ni.ImagesID != nil {
		imgsID = *ni.ImagesID
	}
	i := listUsecase.NewItem{
		ListID:      listID,
		UserID:      claims.Subject,
		Name:        ni.Name,
		Description: ni.Description,
		Address:     ni.Address,
		Point:       np,
		ImagesID:    imgsID,
	}
	res, err := s.core.CreateItem(ctx, i)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrCreateItemBusiness.Error())
		return fmt.Errorf(
			"cannot create item: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	pi := populateItemResponse(res)
	return web.Respond(ctx, w, pi, http.StatusCreated)
}

func (s *Service) UpdateItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	ui := UpdateItem{}
	if err := web.Decode(r, &ui); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrItemUpdateValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	listID, err := getListIDParam(r)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrItemValidateListUUID.Error())
		return err
	}
	itemID, err := getItemIDParam(r)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrItemValidateItemUUID.Error())
		return err
	}
	var up *listUsecase.UpdatePoint
	if ui.Point != nil {
		up = &listUsecase.UpdatePoint{
			Lat: ui.Point.Lat,
			Lng: ui.Point.Lng,
		}
	}
	var imgsID []string
	if ui.ImagesID != nil {
		imgsID = *ui.ImagesID
	}
	i := listUsecase.UpdateItem{
		ID:          itemID,
		ListID:      listID,
		UserID:      claims.Subject,
		Name:        ui.Name,
		Description: ui.Description,
		Address:     ui.Description,
		Point:       up,
		ImagesID:    imgsID,
		Visited:     ui.Visited,
	}
	res, err := s.core.UpdateItemByID(ctx, i)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrUpdateItemBusiness.Error())
		return fmt.Errorf(
			"cannot update item: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	pi := populateItemResponse(res)
	return web.Respond(ctx, w, pi, http.StatusOK)
}

func (s *Service) DeleteItem(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tID := web.GetTraceID(ctx)
	di := struct{}{}
	if err := web.Decode(r, &di); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrItemDeleteValidate.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	itemID, err := getItemIDParam(r)
	if err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrItemValidateItemUUID.Error())
		return err
	}
	if err := s.core.DeleteItemByID(ctx, claims.Subject, itemID); err != nil {
		s.log.Err(err).Str("TraceID", tID).Msg(ErrDeleteItemBusiness.Error())
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
	i := ItemResponse{
		ID:          res.ID,
		ListID:      res.ListID,
		Name:        res.Name,
		Description: res.Description,
		Address:     res.Address,
		Point: PointResponse{
			ID:     res.Point.ID,
			ItemID: res.Point.ItemID,
			Lat:    res.Point.Lat,
			Lng:    res.Point.Lng,
		},
		ImagesID:    &res.ImagesID,
		Visited:     res.Visited,
		DateCreated: res.DateCreated,
		DateUpdated: res.DateUpdated,
	}
	return i
}
