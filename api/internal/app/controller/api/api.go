package api

import (
	"net/http"
	"os"

	authService "github.com/f4mk/api/internal/app/service/auth"
	"github.com/f4mk/api/internal/app/service/user"
	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/internal/pkg/middleware"
	"github.com/f4mk/api/pkg/web"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Config struct {
	Build    string
	Shutdown chan os.Signal
	Log      *zerolog.Logger
	Auth     *auth.Auth
	DB       *sqlx.DB
}

func New(cfg Config) *web.WebApp {

	app := web.New(
		cfg.Shutdown,
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Metrics(),
		middleware.Panics(cfg.Log),
	)

	us := user.NewService(cfg.Log, cfg.DB)
	as := authService.NewService(cfg.Log, cfg.DB, cfg.Auth)

	app.Handle(http.MethodPost, "/user", us.CreateUser)
	app.Handle(http.MethodGet, "/user", us.GetUsers)
	app.Handle(http.MethodGet, "/user/:id", us.GetUser)
	app.Handle(http.MethodPut, "/user/:id", us.UpdateUser, middleware.Authenticate(cfg.Auth))
	app.Handle(http.MethodDelete, "/user/:id", us.DeleteUser, middleware.Authenticate(cfg.Auth))

	app.Handle(http.MethodPost, "/auth/login", as.Login)
	app.Handle(http.MethodPost, "/auth/logout", as.Logout)
	app.Handle(http.MethodPost, "/auth/refresh", as.Refresh)
	app.Handle(http.MethodPost, "/auth/password/reset", as.PasswordReset)

	return app
}
