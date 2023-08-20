package mail

import (
	"context"

	mailSender "github.com/f4mk/api/internal/app/repo/mail"
	"github.com/f4mk/api/internal/app/service/mail"
	"github.com/f4mk/api/pkg/mb"
	"github.com/rs/zerolog"
)

type Config struct {
	Log *zerolog.Logger
	MQ  *mb.Channel
}

type Agent struct {
	service  *mail.Service
	log      *zerolog.Logger
	shutdown context.CancelFunc
}

func New(l *zerolog.Logger, mb *mb.Channel, pb string, pr string) (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r, err := mb.Consume()
	if err != nil {
		cancel()
		return nil, err
	}

	mr := mailSender.NewSender(l, pb, pr)

	ms := mail.NewService(l, mr)

	ma := &Agent{
		service:  ms,
		log:      l,
		shutdown: cancel,
	}

	errCh := make(chan mail.ServeError, 1)
	go ms.Serve(ctx, r, errCh)

	go func() {
		for errMsg := range errCh {
			l.Err(errMsg.Error).Msg("error processing message in mail agent")
		}
	}()
	return ma, nil
}

func (ma *Agent) Shutdown(ctx context.Context) {
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
