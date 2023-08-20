package mail

import (
	"context"
	"fmt"

	"github.com/f4mk/api/internal/app/usecase/mail"
	"github.com/rs/zerolog"
)

type Sender struct {
	log *zerolog.Logger
}

func NewSender(l *zerolog.Logger) *Sender {
	return &Sender{log: l}
}

func (s *Sender) Send(ctx context.Context, l mail.Letter) error {

	fmt.Println(l)
	return nil
}
