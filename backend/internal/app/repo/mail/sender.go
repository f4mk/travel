package mail

import (
	"bytes"
	"embed"
	"html/template"

	mailUsecase "github.com/f4mk/api/internal/app/usecase/mail"
	"github.com/mailjet/mailjet-apiv3-go/v3"
	"github.com/rs/zerolog"
)

//go:embed letter_template.html
var content embed.FS

type Sender struct {
	log        *zerolog.Logger
	publicKey  string
	privateKey string
}

func NewSender(l *zerolog.Logger, pb string, pr string) *Sender {
	return &Sender{log: l, publicKey: pb, privateKey: pr}
}

func (s *Sender) Send(l mailUsecase.Letter) error {

	tmpl, err := template.ParseFS(content, "letter_template.html")
	if err != nil {
		s.log.Err(err).Msg("error parsing letter template from file")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, l); err != nil {
		panic(err)
	}

	mailjetClient := mailjet.NewMailjetClient(s.publicKey, s.privateKey)

	var recipients []mailjet.Recipient
	recipients = append(recipients, mailjet.Recipient{
		Email: l.To,
		Name:  l.Name,
	})

	// TODO: refactor this
	email := &mailjet.InfoSendMail{
		FromEmail:  "noreply@traillyst.com",
		FromName:   "CoolApp",
		Subject:    l.Subject,
		TextPart:   l.Header + "\n" + l.Body + "\n" + l.Link,
		HTMLPart:   buf.String(),
		Recipients: recipients,
	}

	_, err = mailjetClient.SendMail(email)
	if err != nil {
		s.log.Err(err).Msg("error sending email")
		return err
	}

	s.log.Info().Msgf("email sent successfully")

	return nil
}
