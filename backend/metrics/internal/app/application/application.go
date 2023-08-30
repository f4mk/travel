package application

import (
	"net/http"

	"github.com/f4mk/travel/backend/metrics/config"
	app "github.com/f4mk/travel/backend/metrics/internal/app/controller"
	"github.com/f4mk/travel/backend/metrics/internal/app/provider"
	"github.com/f4mk/travel/backend/metrics/internal/app/service"
	metrics "github.com/f4mk/travel/backend/metrics/internal/app/usecase"
	"github.com/f4mk/travel/backend/pkg/utils"
	"github.com/rs/zerolog"
)

func Run(log *zerolog.Logger, cfg *config.Config) error {

	log.Info().Msgf(
		"metrics is listening on: %s",
		utils.GetHost(cfg.API.HostName, cfg.API.Port),
	)

	coll := provider.New(log, utils.GetHost(cfg.Target.HostName, cfg.Target.Port))
	core := metrics.New(log, coll)
	svc := service.New(log, core)

	return http.ListenAndServe(
		utils.GetHost(cfg.API.HostName, cfg.API.Port),
		app.New(app.Config{
			Log:            log,
			MetricsService: svc,
		}),
	)
}
