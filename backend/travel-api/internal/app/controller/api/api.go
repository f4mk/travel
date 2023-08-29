package api

import (
	"net/http"
	"os"
	"time"

	authService "github.com/f4mk/travel/backend/travel-api/internal/app/service/auth"
	userService "github.com/f4mk/travel/backend/travel-api/internal/app/service/user"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/middleware"
	"github.com/f4mk/travel/backend/travel-api/pkg/web"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Shutdown       chan os.Signal
	Log            *zerolog.Logger
	Tracer         trace.Tracer
	Auth           *auth.Auth
	RequestTimeout time.Duration
	RateLimit      int
	AuthService    *authService.Service
	UserService    *userService.Service
}

func New(cfg Config) *web.App {

	app := web.New(
		cfg.Shutdown,
		cfg.RequestTimeout,
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Metrics(),
		middleware.Panics(cfg.Log),
	)

	app.Handle(
		http.MethodPost,
		"/user",
		cfg.UserService.CreateUser,
		middleware.RateLimit(cfg.Log, cfg.RateLimit),
	)
	app.Handle(http.MethodGet, "/user", cfg.UserService.GetUsers)
	app.Handle(http.MethodGet, "/user/:id", cfg.UserService.GetUser)
	app.Handle(
		http.MethodPut,
		"/user/:id",
		cfg.UserService.UpdateUser,
		middleware.RateLimit(cfg.Log, cfg.RateLimit),
		middleware.Authenticate(cfg.Auth),
	)
	app.Handle(
		http.MethodDelete,
		"/user/:id",
		cfg.UserService.DeleteUser,
		middleware.Authenticate(cfg.Auth),
	)

	// TODO: login takes too long, need to do smth
	app.Handle(http.MethodPost, "/auth/login", cfg.AuthService.Login)
	app.Handle(
		http.MethodPost,
		"/auth/logout",
		cfg.AuthService.Logout,
		middleware.Authenticate(cfg.Auth),
	)
	app.Handle(
		http.MethodPost,
		"/auth/logout/all",
		cfg.AuthService.LogoutAll,
		middleware.Authenticate(cfg.Auth),
	)
	app.Handle(
		http.MethodPost,
		"/auth/refresh",
		cfg.AuthService.Refresh,
		middleware.RateLimit(cfg.Log, cfg.RateLimit),
	)
	app.Handle(
		http.MethodPost,
		"/auth/password/change",
		cfg.AuthService.ChangePassword,
		middleware.Authenticate(cfg.Auth),
	)
	app.Handle(
		http.MethodPost,
		"/auth/password/reset",
		cfg.AuthService.PasswordReset,
		middleware.RateLimit(cfg.Log, cfg.RateLimit),
	)
	app.Handle(
		http.MethodPost,
		"/auth/password/reset/submit",
		cfg.AuthService.PasswordResetSubmit,
		middleware.RateLimit(cfg.Log, cfg.RateLimit),
	)

	return app
}
