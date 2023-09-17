package api

import (
	"net/http"

	authService "github.com/f4mk/travel/backend/travel-api/internal/app/service/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/middleware"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

type AuthController struct {
	Log         *zerolog.Logger
	AuthService *authService.Service
	Auth        *auth.Auth
	RateLimit   int
}

func NewAuthController(
	l *zerolog.Logger,
	as *authService.Service,
	a *auth.Auth,
	rl int,
) *AuthController {
	return &AuthController{
		Log:         l,
		AuthService: as,
		Auth:        a,
		RateLimit:   rl,
	}
}

func (ac *AuthController) RegisterRoutes(app *web.App) {
	// TODO: login takes too long
	app.Handle(http.MethodPost, "/auth/login", ac.AuthService.Login)
	app.Handle(
		http.MethodPost,
		"/auth/logout",
		ac.AuthService.Logout,
		middleware.Authenticate(ac.Auth),
	)
	app.Handle(
		http.MethodPost,
		"/auth/logout/all",
		ac.AuthService.LogoutAll,
		middleware.Authenticate(ac.Auth),
	)
	app.Handle(
		http.MethodPost,
		"/auth/refresh",
		ac.AuthService.Refresh,
		middleware.RateLimit(ac.Log, ac.RateLimit),
	)
	app.Handle(
		http.MethodPost,
		"/auth/validate",
		ac.AuthService.Validate,
		middleware.RateLimit(ac.Log, ac.RateLimit),
	)
	app.Handle(
		http.MethodPost,
		"/auth/password/change",
		ac.AuthService.ChangePassword,
		middleware.Authenticate(ac.Auth),
	)
	app.Handle(
		http.MethodPost,
		"/auth/password/reset",
		ac.AuthService.PasswordReset,
		middleware.RateLimit(ac.Log, ac.RateLimit),
	)
	app.Handle(
		http.MethodPost,
		"/auth/password/reset/submit",
		ac.AuthService.PasswordResetSubmit,
		middleware.RateLimit(ac.Log, ac.RateLimit),
	)
}
