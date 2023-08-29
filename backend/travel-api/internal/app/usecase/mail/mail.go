package mail

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Sender interface {
	Send(l Letter) error
}

type Core struct {
	sender Sender
	log    *zerolog.Logger
}

func NewCore(l *zerolog.Logger, s Sender) *Core {
	return &Core{
		sender: s,
		log:    l,
	}
}

func (c *Core) SendMessage(m Message) error {

	sub := "Password reset"
	head := fmt.Sprintf("Hello %s", m.Name)
	body := `You (or somebody on your behalf)
	 have requested a password reset. Is that was not you, just ignore this letter.
	 Otherwise, please, follow the provided link to set a new password.
	 Keep in mind that you cannot request this letter more than once per 10 minutes`

	l := Letter{
		To:      m.Email,
		Name:    m.Name,
		Subject: sub,
		Header:  head,
		Token:   m.ResetToken,
		Body:    body,
	}
	return c.sender.Send(l)
}
