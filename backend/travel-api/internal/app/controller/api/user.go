package api

import (
	"net/http"

	userService "github.com/f4mk/travel/backend/travel-api/internal/app/service/user"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/middleware"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

type UserController struct {
	Log         *zerolog.Logger
	UserService *userService.Service
	Auth        *auth.Auth
	RateLimit   int
}

func (uc *UserController) RegisterRoutes(app *web.App) {
	app.Handle(
		http.MethodPost,
		"/user",
		uc.UserService.CreateUser,
		middleware.RateLimit(uc.Log, uc.RateLimit),
	)
	app.Handle(http.MethodGet, "/user", uc.UserService.GetUsers)
	app.Handle(http.MethodGet, "/user/:id", uc.UserService.GetUser)
	app.Handle(
		http.MethodPut,
		"/user/:id",
		uc.UserService.UpdateUser,
		middleware.RateLimit(uc.Log, uc.RateLimit),
		middleware.Authenticate(uc.Auth),
	)
	app.Handle(
		http.MethodDelete,
		"/user/:id",
		uc.UserService.DeleteUser,
		middleware.Authenticate(uc.Auth),
	)
}
