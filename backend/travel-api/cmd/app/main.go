package main

import (
	"expvar"
	_ "expvar"
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/f4mk/travel/backend/travel-api/config"
	"github.com/f4mk/travel/backend/travel-api/internal/app/application"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/logger"
)

var build = "dev"
var date = "now"
var configPath = "config/.env"

const KEY = "API"

func main() {

	cfg, err := config.New(configPath)
	if err != nil {
		fmt.Println("unable to load config", err)
		os.Exit(1)
	}

	log := logger.New(cfg, KEY)
	log.Info().Msg("Hello world")
	log.Info().Msgf("build: %s", build)
	log.Info().Msgf("date: %s", date)
	log.Info().Msgf("config is loaded from: %s", configPath)

	expvar.NewString("build").Set(build)

	if err := application.Run(build, log, cfg); err != nil {
		log.Err(err).Msgf("unable to launch app: ")
		os.Exit(1)
	}

}
