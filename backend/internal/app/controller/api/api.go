package api

import (
	"net/http"
	"os"
	"time"

	authRepo "github.com/f4mk/api/internal/app/provider/auth"
	userRepo "github.com/f4mk/api/internal/app/provider/user"
	authService "github.com/f4mk/api/internal/app/service/auth"
	userService "github.com/f4mk/api/internal/app/service/user"
	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/internal/pkg/middleware"
	"github.com/f4mk/api/pkg/mb"
	"github.com/f4mk/api/pkg/web"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Config struct {
	Build          string
	Shutdown       chan os.Signal
	Log            *zerolog.Logger
	Auth           *auth.Auth
	DB             *sqlx.DB
	MQ             *mb.Channel
	RequestTimeout time.Duration
	RateLimit      int
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

	ur := userRepo.NewRepo(cfg.Log, cfg.DB)
	ar := authRepo.NewRepo(cfg.Log, cfg.DB)

	us := userService.NewService(cfg.Log, ur)
	as := authService.NewService(cfg.Log, cfg.Auth, ar, cfg.MQ)

	app.Handle(http.MethodPost, "/user", us.CreateUser, middleware.RateLimit(cfg.Log, cfg.RateLimit))
	app.Handle(http.MethodGet, "/user", us.GetUsers)
	app.Handle(http.MethodGet, "/user/:id", us.GetUser)
	app.Handle(
		http.MethodPut,
		"/user/:id",
		us.UpdateUser,
		middleware.RateLimit(cfg.Log, cfg.RateLimit),
		middleware.Authenticate(cfg.Auth),
	)
	app.Handle(http.MethodDelete, "/user/:id", us.DeleteUser, middleware.Authenticate(cfg.Auth))

	// TODO: login takes too long, need to do smth
	app.Handle(http.MethodPost, "/auth/login", as.Login)
	app.Handle(http.MethodPost, "/auth/logout", as.Logout, middleware.Authenticate(cfg.Auth))
	app.Handle(http.MethodPost, "/auth/refresh", as.Refresh, middleware.RateLimit(cfg.Log, cfg.RateLimit))
	app.Handle(http.MethodPost, "/auth/password/change", as.ChangePassword, middleware.Authenticate(cfg.Auth))
	app.Handle(
		http.MethodPost,
		"/auth/password/reset",
		as.PasswordReset,
		middleware.RateLimit(cfg.Log, cfg.RateLimit),
	)
	app.Handle(
		http.MethodPost,
		"/auth/password/reset/submit",
		as.PasswordResetSubmit,
		middleware.RateLimit(cfg.Log, cfg.RateLimit),
	)

	return app
}
