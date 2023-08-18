package application

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/f4mk/api/config"
	"github.com/f4mk/api/internal/app/controller/api"
	"github.com/f4mk/api/internal/app/controller/debug"
	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/internal/pkg/database"
	"github.com/f4mk/api/internal/pkg/keystore"
	"github.com/f4mk/api/internal/pkg/queue"
	"github.com/f4mk/api/pkg/utils"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Run(build string, log *zerolog.Logger, cfg *config.Config) error {

	// -------------------------------------------------------------------------
	// Initializing DB Connection
	log.Info().Msgf(
		"api: initializing database connection %s",
		utils.GetHost(cfg.DB.HostName, cfg.DB.Port),
	)

	db, err := database.Open(database.Config{
		User:        cfg.DB.User,
		Password:    cfg.DB.Password,
		Host:        utils.GetHost(cfg.DB.HostName, cfg.DB.Port),
		Name:        cfg.DB.DBName,
		DisableTLS:  cfg.DB.DisableTLS,
		MaxIdleConn: cfg.DB.MaxIdleConn,
		MaxOpenConn: cfg.DB.MaxOpenConn,
	})

	if err != nil {
		log.Err(err).Msg(ErrInitConnDB.Error())
		return ErrInitConnDB
	}

	defer func() {
		log.Info().Msgf("api: closing db connection %s", utils.GetHost(cfg.DB.HostName, cfg.DB.Port))
		db.Close()
	}()

	// -------------------------------------------------------------------------
	// Initializing message broker connection manager
	cm, err := queue.NewManager(queue.ConnConfig{
		User:     cfg.MessageBroker.User,
		Password: cfg.MessageBroker.Password,
		Host:     utils.GetHost(cfg.MessageBroker.HostName, cfg.MessageBroker.Port),
		Log:      log,
	})
	if err != nil {
		log.Err(err).Msg(ErrCreateBroker.Error())
		return ErrCreateBroker
	}

	mq, err := cm.NewChannel(queue.ChConfig{
		QName:   "resetPasswordLetter",
		WithDLQ: true,
	})
	if err != nil {
		log.Err(err).Msg(ErrCreateQueue.Error())
		return ErrCreateQueue
	}

	// -------------------------------------------------------------------------
	// Creating Auth
	//loading keys
	ks, err := keystore.NewFS(os.DirFS(cfg.Auth.KeyPath))
	if err != nil {
		log.Err(err).Msg(ErrCreateKeyStore.Error())
		return ErrCreateKeyStore
	}

	activeKids := ks.CollectKeyIDs()

	//creating cache
	rdb := redis.NewClient(&redis.Options{
		Addr:         utils.GetHost(cfg.Cache.HostName, cfg.Cache.Port),
		Password:     "",
		DB:           cfg.Cache.DB,
		PoolSize:     cfg.Cache.PoolSize,
		MinIdleConns: cfg.Cache.MinIdleConns,
	})

	pong, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		log.Err(err).Msgf(ErrConnRedis.Error())
		return ErrConnRedis
	}

	log.Info().Msgf("api: connected to redis: %s", pong)

	defer func() {
		log.Info().Msgf("api: closing rdc connection %s", utils.GetHost(cfg.Cache.HostName, cfg.Cache.Port))
		rdb.Close()
	}()

	authCfg := auth.Config{
		ActiveKIDs:      activeKids,
		KeyLookup:       ks,
		Cache:           rdb,
		DB:              db,
		Log:             log,
		AuthDuration:    cfg.Auth.AuthDuration,
		RefreshDuration: cfg.Auth.RefreshDuration,
	}

	auth, err := auth.New(authCfg)

	if err != nil {
		log.Err(err).Msg(ErrConatructAuth.Error())
		return ErrConatructAuth
	}

	// -------------------------------------------------------------------------
	// Start Debug Service
	log.Info().Msgf("debug: initializing debug server: %s", utils.GetHost(cfg.Debug.HostName, cfg.Debug.Port))

	go func() {
		log.Info().Msgf("debug: debug is listening on: %s", utils.GetHost(cfg.Debug.HostName, cfg.Debug.Port))
		if err := http.ListenAndServe(
			utils.GetHost(cfg.Debug.HostName, cfg.Debug.Port),
			debug.New(debug.Config{
				Build: build,
				Log:   log,
				DB:    db,
			}),
		); err != nil {
			log.Err(err).Msgf(ErrRunDebug.Error())
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service
	log.Info().Msgf("api: initializing API server: %s", utils.GetHost(cfg.API.HostName, cfg.API.Port))

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	apiCfg := api.Config{
		Build:          build,
		Shutdown:       shutdown,
		Log:            log,
		Auth:           auth,
		DB:             db,
		MQ:             mq,
		RequestTimeout: cfg.API.RequestTimeout,
		RateLimit:      cfg.API.RateLimit,
	}

	h2s := &http2.Server{}

	api := &http.Server{
		Addr:         utils.GetHost(cfg.API.HostName, cfg.API.Port),
		Handler:      h2c.NewHandler(api.New(apiCfg), h2s),
		ReadTimeout:  cfg.API.ReadTimeout,
		WriteTimeout: cfg.API.WriteTimeout,
	}

	if err := http2.ConfigureServer(api, &http2.Server{}); err != nil {
		log.Err(err).Msg(ErrStartServer.Error())
		return ErrStartServer
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info().Msgf("api: api is listening on: %s", utils.GetHost(cfg.API.HostName, cfg.API.Port))
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown
	select {
	case err := <-serverErrors:
		log.Err(err).Msg(ErrRunServer.Error())
		return ErrRunServer
	case sig := <-shutdown:
		log.Info().Msgf("api: shutting down on signal: %s", sig)
		defer log.Info().Msgf("api: shutdown completed on signal: %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.API.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			log.Err(err).Msg(ErrGracefulShutdown.Error())
			return ErrGracefulShutdown
		}
	}

	return nil
}
