package mail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	mailUsecase "github.com/f4mk/api/internal/app/usecase/mail"
	"github.com/f4mk/api/pkg/mb"
	"github.com/rs/zerolog"
)

type Service struct {
	core   *mailUsecase.Core
	log    *zerolog.Logger
	doneCh chan struct{}
}

func NewService(l *zerolog.Logger, s mailUsecase.Sender) *Service {
	return &Service{
		core:   mailUsecase.NewCore(l, s),
		log:    l,
		doneCh: make(chan struct{}),
	}
}

func (s *Service) Serve(ctx context.Context, rx <-chan mb.Message, errCh chan ServeError) {
	for {
		select {
		case msg, ok := <-rx:
			if !ok {
				close(s.doneCh)
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
				// send error outside
				err = sendError(errCh, fmt.Errorf(ErrNackReqMessage.Error(), err), msg.Body)
				if err != nil {
					s.log.Warn().Msg(ErrChanFull.Error())
				}
			}

			// process the letter
			if err := s.core.SendMessage(m); err != nil {
				s.log.Err(err).Msg(ErrUsecaseLayer.Error())

				// TODO: handle retries here

				// if reached here, it means the retry limit is exceeded
				if err := msg.Nack(false, false); err != nil {
					s.log.Err(err).Msg(ErrNackMessage.Error())
				}
				err = sendError(errCh, fmt.Errorf(ErrNackReqMessage.Error(), err), msg.Body)
				if err != nil {
					s.log.Warn().Msg(ErrChanFull.Error())
				}
			}
			if err := msg.Ack(false); err != nil {
				// TODO: retry? nothing really can do here, the message was already processed
				s.log.Err(err).Msg(ErrAckMessage.Error())
			}

		case <-ctx.Done():
			close(s.doneCh)
			return
		}
	}
}

func (s *Service) Stop() <-chan struct{} {
	return s.doneCh
}

func sendError(errCh chan<- ServeError, errMsg error, payload []byte) error {
	select {
	case errCh <- ServeError{
		Error:   errMsg,
		Payload: payload,
	}:
		return nil
	default:
		return errors.New("message error channel is full, might be an error")
	}
}

// func parseCountFromHeader(msg mb.Message) (int32, error) {
// 	if h, ok := msg.Headers["x-death"]; ok {
// 		v, ok := h.([]interface{})
// 		if !ok {
// 			return 0, fmt.Errorf("couldnt parse slice: %w", ErrHeaderFormat)

// 		}
// 		if len(v) == 0 {
// 			return 0, fmt.Errorf("slice is empty: %w", ErrHeaderFormat)
// 		}

// 		details, ok := v[0].(amqp.Table)
// 		if !ok {
// 			return 0, fmt.Errorf("couldnt parse ampq table: %w", ErrHeaderFormat)
// 		}

// 		count, ok := details["count"].(int32)
// 		if !ok {
// 			return 0, ErrGetCount
// 		}
// 		return count, nil
// 	}
// 	return 0, ErrMissingHeader
// }
