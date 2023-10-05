package mail

import (
	"context"

	"github.com/f4mk/travel/backend/travel-api/internal/app/service/mail"
	"github.com/rs/zerolog"
)

type Config struct {
	Log         *zerolog.Logger
	MailService *mail.Service
}

type Agent struct {
	service  *mail.Service
	log      *zerolog.Logger
	shutdown context.CancelFunc
}

func New(cfg Config) (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())

	ma := Agent{
		service:  cfg.MailService,
		log:      cfg.Log,
		shutdown: cancel,
	}

	errMsgCh := make(chan mail.ServeError, 1)
	errServiceCh := make(chan error)
	// TODO: needs to reconnect
	go ma.service.Serve(ctx, errMsgCh, errServiceCh)

	go func() {
		for errMsg := range errMsgCh {
			ma.log.Err(errMsg.Error).Msg("error processing message in mail agent")
		}
	}()
	err, ok := <-errServiceCh
	if ok {
		if err != nil {
			cancel()
			return nil, err
		}
	}
	return &ma, nil
}

func (ma *Agent) Shutdown(ctx context.Context) {
	ma.log.Warn().Msg("shutting down mail agent")
	ma.shutdown()
	doneCh := make(chan struct{})
	go func() {
		<-ma.service.Stop()
		close(doneCh)
	}()

	select {
	case <-doneCh:
	case <-ctx.Done():
		ma.log.Error().Msg("error graceful shutdown mail service")
	}
}
