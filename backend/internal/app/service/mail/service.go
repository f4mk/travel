package mail

import (
	"context"
	"encoding/json"
	"fmt"

	mailUsecase "github.com/f4mk/api/internal/app/usecase/mail"
	"github.com/f4mk/api/pkg/mb"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

const numberOfRetries = 3

type Service struct {
	core       *mailUsecase.Core
	log        *zerolog.Logger
	shutdownCh chan struct{}
}

func NewService(l *zerolog.Logger, s mailUsecase.Sender) *Service {

	return &Service{
		core:       mailUsecase.NewCore(l, s),
		log:        l,
		shutdownCh: make(chan struct{}),
	}

}

func (s *Service) Serve(ctx context.Context, rx <-chan mb.Message, errCh chan ServeError) {
	defer close(s.shutdownCh)
	for {
		select {
		case msg, ok := <-rx:
			if !ok {
				return
			}
			m := mailUsecase.Message{}

			err := json.Unmarshal(msg.Body, &m)
			if err != nil {
				s.log.Err(err).Msg(ErrParseMessage.Error())
				// makes no sense to requeue due to invalid json
				if err := msg.Nack(false, false); err != nil {
					s.log.Err(err).Msg(ErrNackReqMessage.Error())
				}
				errCh <- ServeError{
					Error:   fmt.Errorf(ErrNackReqMessage.Error(), err),
					Payload: msg.Body,
				}
				continue
			}

			if err := s.core.SendMessage(ctx, m); err != nil {
				c, err := parseCountFromHeader(msg)
				if err != nil || c < numberOfRetries {
					if err := msg.Nack(false, true); err != nil {
						s.log.Err(err).Msg(ErrNackReqMessage.Error())
					}
					continue
				}
				// if reached here, it means the retry limit is exceeded
				if err := msg.Nack(false, false); err != nil {
					s.log.Err(err).Msg(ErrNackMessage.Error())
				}
				continue
			}
			if err := msg.Ack(false); err != nil {
				s.log.Err(err).Msg(ErrAckMessage.Error())
			}

		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) Stop() <-chan struct{} {
	return s.shutdownCh
}

func parseCountFromHeader(msg mb.Message) (int32, error) {
	if h, ok := msg.Headers["x-death"]; ok {
		v, ok := h.([]interface{})
		if !ok {
			return 0, fmt.Errorf("couldnt parse slice: %w", ErrHeaderFormat)

		}
		if len(v) == 0 {
			return 0, fmt.Errorf("slice is empty: %w", ErrHeaderFormat)
		}

		details, ok := v[0].(amqp.Table)
		if !ok {
			return 0, fmt.Errorf("couldnt parse ampq table: %w", ErrHeaderFormat)
		}

		count, ok := details["count"].(int32)
		if !ok {
			return 0, ErrGetCount
		}
		return count, nil
	}
	return 0, ErrMissingHeader
}
