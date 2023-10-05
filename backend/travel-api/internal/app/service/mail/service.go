package mail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	queue "github.com/f4mk/travel/backend/pkg/mb"
	mailUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/mail"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/messages"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
)

type Service struct {
	core   *mailUsecase.Core
	log    *zerolog.Logger
	doneCh chan struct{}
	mq     *queue.Channel
}

func NewService(l *zerolog.Logger, c *mailUsecase.Core, mq *queue.Channel) *Service {
	return &Service{
		core:   c,
		log:    l,
		doneCh: make(chan struct{}),
		mq:     mq,
	}
}

func (s *Service) Serve(ctx context.Context, errMsgCh chan<- ServeError, errServiceCh chan<- error) {
	rx, err := s.mq.Consume()
	if err != nil {
		errServiceCh <- err
		return
	}
	// init complete
	close(errServiceCh)

	for {
		select {
		case msg, ok := <-rx:
			if !ok {
				close(s.doneCh)
				s.log.Warn().Msg("rx channel got closed, returning from service")
				return
			}

			m := messages.Message{}
			err := json.Unmarshal(msg.Body, &m)
			tID := m.ID
			ctx, span := web.AddSpan(ctx, "service.mail.serve", attribute.String("TraceID", tID))
			defer span.End()

			v := web.Values{
				TraceID: tID,
				Now:     time.Now().UTC(),
			}
			ctx = web.SetValues(ctx, &v)

			if err != nil {
				s.log.Err(err).Str("TraceID", tID).Msgf("error decoding message: %s", ErrParseMessage.Error())
				// makes no sense to requeue due to invalid json
				if err := msg.Nack(false, false); err != nil {
					s.log.Err(err).Str("TraceID", tID).Msgf("error NACKing after decoding: %s", ErrNackMessage.Error())
				}
				// send error outside
				err = sendError(errMsgCh, fmt.Errorf(ErrNackReqMessage.Error(), err), msg.Body)
				if err != nil {
					s.log.Warn().Str("TraceID", tID).Msgf("error sending err after decoding: %s", ErrChanFull.Error())
				}
			}

			if m.Type == messages.ResetPassword {
				mReset := mailUsecase.MessageReset{
					Email:      strings.ToLower(m.Email),
					Name:       m.Name,
					ResetToken: m.Token,
				}
				err = s.core.SendResetMessage(ctx, mReset)
			} else if m.Type == messages.RegisterVerify {
				mVerify := mailUsecase.MessageVerify{
					Email:       strings.ToLower(m.Email),
					Name:        m.Name,
					VerifyToken: m.Token,
				}
				err = s.core.SendVerifyMessage(ctx, mVerify)
			}
			// process the letter
			if err != nil {
				s.log.Err(err).Str("TraceID", tID).Msg(ErrUsecaseLayer.Error())
				// TODO: handle retries here
				// if reached here, it means the retry limit is exceeded
				if err := msg.Nack(false, false); err != nil {
					s.log.Err(err).Str("TraceID", tID).Msg(ErrNackMessage.Error())
				}
				err = sendError(errMsgCh, fmt.Errorf(ErrNackReqMessage.Error(), err), msg.Body)
				if err != nil {
					s.log.Warn().Str("TraceID", tID).Msg(ErrChanFull.Error())
				}
			} else if err := msg.Ack(false); err != nil {
				// TODO: retry? nothing really can do here, the message was already processed
				s.log.Err(err).Str("TraceID", tID).Msg(ErrAckMessage.Error())
			}

		case <-ctx.Done():
			s.log.Warn().Msg("shutting down mail service due to ctx done")
			close(s.doneCh)
			return
		}
	}
}

func (s *Service) Stop() <-chan struct{} {
	s.log.Warn().Msg("stopping mail service")
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
