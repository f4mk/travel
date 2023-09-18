package image

import (
	"context"
	"fmt"
	"net/http"

	imageUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/image"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	authPkg "github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"

	"github.com/rs/zerolog"
)

type Service struct {
	log  *zerolog.Logger
	auth *authPkg.Auth
	core *imageUsecase.Core
}

func NewService(
	l *zerolog.Logger,
	a *authPkg.Auth,
	c *imageUsecase.Core,
) *Service {

	return &Service{
		log:  l,
		auth: a,
		core: c,
	}
}

func (s *Service) Serve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fileID := web.Param(r, "fname")
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	res, err := s.core.GetImageByID(ctx, fileID, claims.Subject)
	if err != nil {
		s.log.Err(err).Msg(ErrGetImageBusiness.Error())
		return fmt.Errorf(
			"cannot get image: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}

	return web.RespondRaw(ctx, w, res, http.StatusOK, "image/webp")
}

func (s *Service) Store(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	listID := web.Param(r, "listID")
	itemID := web.Param(r, "itemID")
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	res, err := s.core.StoreImage(ctx, itemID, listID, claims.Subject)
	if err != nil {
		s.log.Err(err).Msg(ErrPostImageBusiness.Error())
		return fmt.Errorf(
			"cannot store image: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}

	return web.Respond(ctx, w, res, http.StatusCreated)
}
