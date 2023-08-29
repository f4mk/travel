package application

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/f4mk/travel/backend/pkg/mb"
	"github.com/f4mk/travel/backend/pkg/utils"
	"github.com/f4mk/travel/backend/travel-api/config"
	"github.com/f4mk/travel/backend/travel-api/internal/app/controller/api"
	"github.com/f4mk/travel/backend/travel-api/internal/app/controller/debug"
	"github.com/f4mk/travel/backend/travel-api/internal/app/controller/mail"
	authRepo "github.com/f4mk/travel/backend/travel-api/internal/app/provider/auth"
	mailSender "github.com/f4mk/travel/backend/travel-api/internal/app/provider/mail"
	userRepo "github.com/f4mk/travel/backend/travel-api/internal/app/provider/user"
	authService "github.com/f4mk/travel/backend/travel-api/internal/app/service/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/app/service/check"
	mailService "github.com/f4mk/travel/backend/travel-api/internal/app/service/mail"
	userService "github.com/f4mk/travel/backend/travel-api/internal/app/service/user"
	authUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/auth"
	userUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/user"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/keystore"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/tracer"
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
	log.Info().Msg("api: initializing mb manager")
	cm, err := mb.NewManager(mb.ConnConfig{
		User:     cfg.MessageBroker.User,
		Password: cfg.MessageBroker.Password,
		Host:     utils.GetHost(cfg.MessageBroker.HostName, cfg.MessageBroker.Port),
		Log:      log,
	})
	if err != nil {
		log.Err(err).Msg(ErrCreateBroker.Error())
		return ErrCreateBroker
	}
	defer cm.Close()

	mq, err := cm.NewChannel(mb.ChConfig{
		QName:   "resetPasswordLetter",
		WithDLQ: true,
	})
	if err != nil {
		log.Err(err).Msg(ErrCreateQueue.Error())
		return ErrCreateQueue
	}
	defer mq.Close()
	// -------------------------------------------------------------------------
	// Starting Mail service

	mr := mailSender.NewSender(
		log,
		cfg.MailService.PublicKey,
		cfg.MailService.PrivateKey,
		cfg.Service.DomainName,
	)

	ms := mailService.NewService(log, mr, mq)

	mailAgent, err := mail.New(mail.Config{
		Log:         log,
		MailService: ms,
	})
	if err != nil {
		log.Err(err).Msg(ErrCreateMailServer.Error())
		return ErrCreateMailServer
	}

	// -------------------------------------------------------------------------
	// Creating Auth
	// loading keys
	ks, err := keystore.NewFS(os.DirFS(cfg.Auth.KeyPath))
	if err != nil {
		log.Err(err).Msg(ErrCreateKeyStore.Error())
		return ErrCreateKeyStore
	}

	activeKids := ks.CollectKeyIDs()

	// creating cache
	redis := redis.NewClient(&redis.Options{
		Addr:         utils.GetHost(cfg.Cache.HostName, cfg.Cache.Port),
		Password:     "",
		DB:           cfg.Cache.DB,
		PoolSize:     cfg.Cache.PoolSize,
		MinIdleConns: cfg.Cache.MinIdleConns,
	})

	pong, err := redis.Ping(context.TODO()).Result()
	if err != nil {
		log.Err(err).Msgf(ErrConnRedis.Error())
		return ErrConnRedis
	}

	log.Info().Msgf("api: connected to redis: %s", pong)

	defer func() {
		log.Info().Msgf("api: closing rdc connection %s", utils.GetHost(cfg.Cache.HostName, cfg.Cache.Port))
		redis.Close()
	}()

	authCfg := auth.Config{
		ActiveKIDs:      activeKids,
		KeyLookup:       ks,
		Cache:           redis,
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
	// Start Tracing Support
	log.Info().Msg("api: initializing otel tracing")

	traceProvider, err := tracer.NewProvider(
		cfg.Service.ServiceName,
		utils.GetHost(cfg.Telemetry.HostName, cfg.Telemetry.Port),
		cfg.Telemetry.Prob,
	)
	if err != nil {
		log.Err(err).Msg("starting tracing")
		return err
	}
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Err(err).Msg("error shutting down otel trace provider")
		}
	}()

	// TODO: add spans across app
	tracer := traceProvider.Tracer("api")

	// -------------------------------------------------------------------------
	// Start Debug Service
	log.Info().Msgf("debug: initializing debug server: %s", utils.GetHost(cfg.Debug.HostName, cfg.Debug.Port))
	check := check.NewService(build, log, db)
	debugErrors := make(chan error, 1)

	go func() {
		log.Info().Msgf("debug: debug is listening on: %s", utils.GetHost(cfg.Debug.HostName, cfg.Debug.Port))
		debugErrors <- http.ListenAndServe(
			utils.GetHost(cfg.Debug.HostName, cfg.Debug.Port),
			debug.New(debug.Config{
				Build:   build,
				Log:     log,
				DB:      db,
				Service: check,
			}),
		)
	}()

	// -------------------------------------------------------------------------
	// Start API Service
	log.Info().Msgf("api: initializing API server: %s", utils.GetHost(cfg.API.HostName, cfg.API.Port))

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	ur := userRepo.NewRepo(log, db)
	ar := authRepo.NewRepo(log, db)

	uc := userUsecase.NewCore(log, ur)
	us := userService.NewService(log, uc)

	ac := authUsecase.NewCore(log, ar)
	as := authService.NewService(log, auth, ac, mq)

	apiCfg := api.Config{
		Shutdown:       shutdown,
		Log:            log,
		Tracer:         tracer,
		Auth:           auth,
		RequestTimeout: cfg.API.RequestTimeout,
		RateLimit:      cfg.API.RateLimit,
		AuthService:    as,
		UserService:    us,
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
	case err := <-debugErrors:
		log.Err(err).Msg(ErrRunDebug.Error())
		return ErrRunDebug
	case sig := <-shutdown:
		log.Info().Msgf("api: shutting down on signal: %s", sig)
		defer log.Info().Msgf("api: shutdown completed on signal: %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.API.ShutdownTimeout)
		defer cancel()

		mailAgent.Shutdown(ctx)

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			log.Err(err).Msg(ErrGracefulShutdown.Error())
			return ErrGracefulShutdown
		}
	}

	return nil
}
