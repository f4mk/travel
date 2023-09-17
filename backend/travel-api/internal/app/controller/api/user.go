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

func NewUserController(
	l *zerolog.Logger,
	us *userService.Service,
	a *auth.Auth,
	rl int,
) *UserController {
	return &UserController{
		Log:         l,
		UserService: us,
		Auth:        a,
		RateLimit:   rl,
	}
}

func (uc *UserController) RegisterRoutes(app *web.App) {
	app.Handle(
		http.MethodPost,
		"/users",
		uc.UserService.CreateUser,
		middleware.RateLimit(uc.Log, uc.RateLimit),
	)
	app.Handle(http.MethodGet, "/users/me", uc.UserService.GetMe, middleware.Authenticate(uc.Auth))
	app.Handle(http.MethodGet, "/users/:id", uc.UserService.GetUser, middleware.Authenticate(uc.Auth))
	app.Handle(
		http.MethodPut,
		"/users",
		uc.UserService.UpdateUser,
		middleware.RateLimit(uc.Log, uc.RateLimit),
		middleware.Authenticate(uc.Auth),
	)
	app.Handle(
		http.MethodDelete,
		"/users",
		uc.UserService.DeleteUser,
		middleware.Authenticate(uc.Auth),
	)
	app.Handle(
		http.MethodPost,
		"/users/verify",
		uc.UserService.VerifyUser,
		middleware.RateLimit(uc.Log, uc.RateLimit),
	)
}
