package application

import (
	"context"
	"crypto/tls"
	"fmt"
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
	"github.com/rs/zerolog"
	"golang.org/x/net/http2"
)

func Run(build string, log *zerolog.Logger, cfg *config.Config) error {

	// -------------------------------------------------------------------------
	// Creating Auth
	ks, err := keystore.NewFS(os.DirFS(cfg.Auth.KeyPath))
	if err != nil {
		return fmt.Errorf("api: error creating keystore: %w", err)
	}

	activeKids := ks.CollectKeyIDs()

	auth, err := auth.New(activeKids, ks)

	if err != nil {
		return fmt.Errorf("api: error constructing auth: %w", err)
	}

	// -------------------------------------------------------------------------
	// Initializing DB Connection
	log.Info().Msgf(
		"api: initializing database connection %s",
		getHost(cfg.DB.HostName, cfg.DB.Port),
	)

	db, err := database.Open(database.Config{
		User:        cfg.DB.User,
		Password:    cfg.DB.Password,
		Host:        getHost(cfg.DB.HostName, cfg.DB.Port),
		Name:        cfg.DB.DBName,
		DisableTLS:  cfg.DB.DisableTLS,
		MaxIdleConn: cfg.DB.MaxIdleConn,
		MaxOpenConn: cfg.DB.MaxOpenConn,
	})

	if err != nil {
		return fmt.Errorf("api: error initializing connection to db: %w", err)
	}

	defer func() {
		log.Info().Msgf("api: closing db connection %s", getHost(cfg.DB.HostName, cfg.DB.Port))
		db.Close()
	}()

	// -------------------------------------------------------------------------
	// Start Debug Service
	log.Info().Msgf("debug: initializing debug server: %s", getHost(cfg.Debug.HostName, cfg.Debug.Port))

	go func() {
		log.Info().Msgf("debug: debug is listening on: %s", getHost(cfg.Debug.HostName, cfg.Debug.Port))
		if err := http.ListenAndServe(
			getHost(cfg.Debug.HostName, cfg.Debug.Port),
			debug.New(debug.Config{
				Build: build,
				Log:   log,
				DB:    db,
			}),
		); err != nil {
			log.Err(err).Msgf("debug: error debug server on: %s", getHost(cfg.Debug.HostName, cfg.Debug.Port))
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service
	log.Info().Msgf("api: initializing API server: %s", getHost(cfg.Api.HostName, cfg.Api.Port))

	serverTLSCert, err := tls.LoadX509KeyPair(cfg.Api.CertFile, cfg.Api.KeyFile)

	if err != nil {
		log.Err(err).Msg("Error loading certificate and key file:")
		return err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	apiCfg := api.Config{
		Build:    build,
		Shutdown: shutdown,
		Log:      log,
		Auth:     auth,
		DB:       db,
	}

	api := &http.Server{
		Addr:      getHost(cfg.Api.HostName, cfg.Api.Port),
		Handler:   api.New(apiCfg),
		TLSConfig: tlsConfig,
	}

	if err := http2.ConfigureServer(api, &http2.Server{}); err != nil {
		log.Err(err).Msg("Error starting http2 server:")
		return err
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info().Msgf("api: api is listening on: %s", getHost(cfg.Api.HostName, cfg.Api.Port))
		serverErrors <- api.ListenAndServeTLS("", "")
	}()

	// -------------------------------------------------------------------------
	// Shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("api: server error: %w", err)
	case sig := <-shutdown:
		log.Info().Msgf("api: shutting down on signal: %s", sig)
		defer log.Info().Msgf("api: shutdown completed on signal: %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Api.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("api: could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func getHost(hostName string, port string) string {

	return hostName + ":" + port
}
