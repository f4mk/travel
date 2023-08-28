package main

import (
	"expvar"
	_ "expvar"
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/f4mk/api/cmd/app/metrics/application"
	"github.com/f4mk/api/config"
	"github.com/f4mk/api/internal/pkg/logger"
)

var configPath = "config/.env"
var build = "dev"
var date = "now"

const KEY = "METRICS"

func main() {

	// TODO: add logging
	cfg, err := config.New(configPath)
	if err != nil {
		fmt.Println("unable to load config", err)
		os.Exit(1)
	}

	expvar.NewString("build").Set(build)

	log := logger.New(cfg, KEY)
	log.Info().Msg("Hello metrics")
	log.Info().Msgf("build: %s", build)
	log.Info().Msgf("date: %s", date)
	log.Info().Msgf("config is loaded from: %s", configPath)

	if err := application.Run(log, cfg); err != nil {
		log.Err(err).Msgf("metrics service shutdown")
		os.Exit(1)
	}
}
