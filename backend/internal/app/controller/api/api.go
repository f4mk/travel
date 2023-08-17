package api

import (
	"net/http"
	"os"
	"time"

	authRepo "github.com/f4mk/api/internal/app/repo/auth"
	userRepo "github.com/f4mk/api/internal/app/repo/user"
	authService "github.com/f4mk/api/internal/app/service/auth"
	userService "github.com/f4mk/api/internal/app/service/user"
	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/internal/pkg/middleware"
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
	as := authService.NewService(cfg.Log, cfg.Auth, ar)

	app.Handle(http.MethodPost, "/user", us.CreateUser, middleware.RateLimit(cfg.Log, cfg.RateLimit))
	app.Handle(http.MethodGet, "/user", us.GetUsers, middleware.RateLimit(cfg.Log, cfg.RateLimit))
	app.Handle(http.MethodGet, "/user/:id", us.GetUser)
	app.Handle(
		http.MethodPut,
		"/user/:id",
		us.UpdateUser,
		middleware.RateLimit(cfg.Log, cfg.RateLimit),
		middleware.Authenticate(cfg.Auth),
	)
	app.Handle(http.MethodDelete, "/user/:id", us.DeleteUser, middleware.Authenticate(cfg.Auth))

	app.Handle(http.MethodPost, "/auth/login", as.Login)
	app.Handle(http.MethodPost, "/auth/logout", as.Logout, middleware.Authenticate(cfg.Auth))
	app.Handle(http.MethodPost, "/auth/changepass", as.ChangePassword, middleware.Authenticate(cfg.Auth))
	app.Handle(http.MethodPost, "/auth/refresh", as.Refresh, middleware.RateLimit(cfg.Log, cfg.RateLimit))
	// TODO:
	app.Handle(
		http.MethodPost,
		"/auth/password/reset",
		as.PasswordReset,
		middleware.RateLimit(cfg.Log, cfg.RateLimit),
	)

	return app
}
