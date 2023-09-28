package image

import (
	"context"
	"fmt"
	"io"
	"net/http"

	imageUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/image"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	authPkg "github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"

	"github.com/rs/zerolog"
)

const meg16 = 16 << 20 //16Mib

type Service struct {
	log  *zerolog.Logger
	auth *authPkg.Auth
	core *imageUsecase.Core
	sem  chan struct{}
}

func NewService(
	l *zerolog.Logger,
	a *authPkg.Auth,
	c *imageUsecase.Core,
	m int16,
) *Service {
	return &Service{
		log:  l,
		auth: a,
		core: c,
		sem:  make(chan struct{}, m),
	}
}

func (s *Service) Serve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fileID := web.Param(r, "fname")
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	reader, err := s.core.GetImageByID(ctx, fileID, claims.Subject)
	if err != nil {
		s.log.Err(err).Msg(ErrGetImageBusiness.Error())
		return fmt.Errorf(
			"cannot get image: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}
	defer reader.Close()

	return web.RespondRaw(ctx, w, reader, http.StatusOK, "image/webp")
}

func (s *Service) Store(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	listID := web.Param(r, "listID")
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	err = r.ParseMultipartForm(meg16)
	if err != nil {
		s.log.Err(err).Msg(ErrPostImageDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	files := r.MultipartForm.File["images"]
	if len(files) == 0 || len(files) > 5 {
		s.log.Error().Msg(ErrPostImageDecodeLen.Error())
		return web.NewRequestError(
			ErrPostImageDecodeLen,
			http.StatusBadRequest,
		)
	}

	var imageStreams []io.Reader

	for _, fileHeader := range files {
		s.sem <- struct{}{}
		file, err := fileHeader.Open()
		if err != nil {
			<-s.sem
			s.log.Err(err).Msg(ErrPostImageRead.Error())
			return web.NewRequestError(
				err,
				http.StatusBadRequest,
			)
		}
		defer func() {
			file.Close()
			<-s.sem
		}()

		imageStreams = append(imageStreams, file)
	}

	res, err := s.core.StoreImages(ctx, imageStreams, listID, claims.Subject)
	if err != nil {
		s.log.Err(err).Msg(ErrPostImageBusiness.Error())
		return fmt.Errorf(
			"cannot store image: %w",
			web.GetResponseErrorFromBusiness(err),
		)
	}

	return web.Respond(ctx, w, res, http.StatusCreated)
}
