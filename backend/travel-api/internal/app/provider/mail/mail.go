package mail

import (
	"bytes"
	"embed"
	"html/template"
	"net/url"

	mailUsecase "github.com/f4mk/travel/backend/travel-api/internal/app/usecase/mail"
	"github.com/mailjet/mailjet-apiv3-go/v3"
	"github.com/rs/zerolog"
)

//go:embed letter_template.html
var content embed.FS

type Sender struct {
	log        *zerolog.Logger
	dName      string
	publicKey  string
	privateKey string
}

func NewSender(l *zerolog.Logger, pb string, pr string, dn string) *Sender {
	return &Sender{log: l, dName: dn, publicKey: pb, privateKey: pr}
}

func (s *Sender) Send(l mailUsecase.Letter) error {

	tmpl, err := template.ParseFS(content, "letter_template.html")
	if err != nil {
		s.log.Err(err).Msg("error parsing letter template from file")
	}

	q := make(url.Values)
	q.Set("token", l.Token)

	u := &url.URL{
		Scheme:   "https",
		Host:     s.dName,
		Path:     "/password/reset",
		RawQuery: q.Encode(),
	}

	link := u.String()
	letter := struct {
		To      string
		Name    string
		Subject string
		Header  string
		Body    string
		Link    string
		Domain  string
	}{
		To:      l.To,
		Name:    l.Name,
		Subject: l.Subject,
		Header:  l.Header,
		Body:    l.Body,
		Link:    link,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, letter); err != nil {
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
		TextPart:   l.Header + "\n" + l.Body + "\n" + link,
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
