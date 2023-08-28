package application

import (
	"net/http"

	app "github.com/f4mk/api/cmd/app/metrics/controller"
	"github.com/f4mk/api/cmd/app/metrics/provider"
	"github.com/f4mk/api/cmd/app/metrics/service"
	metrics "github.com/f4mk/api/cmd/app/metrics/usecase"
	"github.com/f4mk/api/config"
	"github.com/rs/zerolog"
)

func Run(log *zerolog.Logger, cfg *config.Config) error {

	log.Info().Msgf(
		"metrics is listening on: %s",
		// TODO: use personal .env
		// utils.GetHost(cfg.Metrics.HostName, cfg.Metrics.Port),
		"0.0.0.0:8091",
	)

	// TODO: use personal .env
	// coll := provider.New(log, utils.GetHost(cfg.Debug.HostName, cfg.Debug.Port))
	coll := provider.New(log, "travel-api:8081")
	core := metrics.New(log, coll)
	svc := service.New(log, core)

	return http.ListenAndServe(
		// TODO: use personal .env
		// utils.GetHost(cfg.Metrics.HostName, cfg.Metrics.Port),
		"0.0.0.0:8091",
		app.New(app.Config{
			Log:            log,
			MetricsService: svc,
		}),
	)
}
