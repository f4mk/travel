package api

import (
	"net/http"

	imageService "github.com/f4mk/travel/backend/travel-api/internal/app/service/image"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/middleware"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

type ImageController struct {
	Log          *zerolog.Logger
	ImageService *imageService.Service
	Auth         *auth.Auth
	RateLimit    int
}

func NewImageController(
	l *zerolog.Logger,
	is *imageService.Service,
	a *auth.Auth,
	rl int,
) *ImageController {
	return &ImageController{
		Log:          l,
		ImageService: is,
		Auth:         a,
		RateLimit:    rl,
	}
}

func (ic *ImageController) RegisterRoutes(app *web.App) {
	app.Handle(http.MethodGet, "/images/:fname", ic.ImageService.Serve, middleware.Authenticate(ic.Auth))
	app.Handle(http.MethodPost, "/images/upload/:listID", ic.ImageService.Store, middleware.Authenticate(ic.Auth))
}
