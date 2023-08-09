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

	user := user.NewService(cfg.Log, cfg.DB)
	authService := authService.NewService(cfg.Log, cfg.DB, cfg.Auth)

	app.Handle(http.MethodPost, "/user", user.CreateUser)
	app.Handle(http.MethodGet, "/user", user.GetUsers)
	app.Handle(http.MethodGet, "/user/:id", user.GetUser)
	app.Handle(http.MethodPut, "/user/:id", user.UpdateUser, middleware.Authenticate(cfg.Auth))
	app.Handle(http.MethodDelete, "/user/:id", user.DeleteUser, middleware.Authenticate(cfg.Auth))

	app.Handle(http.MethodPost, "/auth/login", authService.Login)
	app.Handle(http.MethodPost, "/auth/logout", authService.Logout)
	app.Handle(http.MethodPost, "/auth/refresh", authService.Refresh)
	app.Handle(http.MethodPost, "/auth/password/reset", authService.PasswordReset)

	return app
}
