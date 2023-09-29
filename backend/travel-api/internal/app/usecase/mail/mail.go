package mail

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type Sender interface {
	SendResetPwdEmail(ctx context.Context, l Letter) error
	SendRegisterEmail(ctx context.Context, l Letter) error
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

func (c *Core) SendResetMessage(ctx context.Context, m MessageReset) error {

	sub := "Password reset"
	head := fmt.Sprintf("Hello %s", m.Name)
	body := `You (or somebody on your behalf)
	 have requested a password reset. Is that was not you, just ignore this letter.
	 Otherwise, please, follow the provided link to set a new password.`

	l := Letter{
		To:      m.Email,
		Name:    m.Name,
		Subject: sub,
		Header:  head,
		Token:   m.ResetToken,
		Body:    body,
	}
	return c.sender.SendResetPwdEmail(ctx, l)
}

func (c *Core) SendVerifyMessage(ctx context.Context, m MessageVerify) error {

	sub := "Account created"
	head := fmt.Sprintf("Hello %s", m.Name)
	body := `You (or somebody on your behalf)
	 have registered on Traillyst. This is your verification letter.
	 Please, follow the provided link in order to verife your account.`

	l := Letter{
		To:      m.Email,
		Name:    m.Name,
		Subject: sub,
		Header:  head,
		Token:   m.VerifyToken,
		Body:    body,
	}
	return c.sender.SendRegisterEmail(ctx, l)
}
