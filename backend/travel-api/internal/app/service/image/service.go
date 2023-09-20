package image

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
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
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		s.log.Err(err).Msgf(auth.ErrGetClaims.Error())
		return auth.ErrGetClaims
	}
	err = r.ParseMultipartForm(16 << 20) // Limit: 16MB
	if err != nil {
		s.log.Err(err).Msg(ErrPostImageDecode.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		s.log.Info().Msg(ErrPostImageDecodeLen.Error())
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	var imageStreams []io.Reader

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			s.log.Err(err).Msg(ErrPostImageRead.Error())
			return web.NewRequestError(
				err,
				http.StatusBadRequest,
			)
		}
		defer file.Close()

		stream, err := streamFile(file)
		if err != nil {
			s.log.Err(err).Msg(ErrPostImageReadContent.Error())
			return web.NewRequestError(
				err,
				http.StatusBadRequest,
			)
		}
		imageStreams = append(imageStreams, stream)
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

func streamFile(file multipart.File) (io.Reader, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
