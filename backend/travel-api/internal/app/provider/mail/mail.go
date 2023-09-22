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
var letterTmpl embed.FS

type Sender struct {
	log   *zerolog.Logger
	dName string
	mail  *mailjet.Client
}

func NewSender(l *zerolog.Logger, m *mailjet.Client, dn string) *Sender {
	return &Sender{log: l, mail: m, dName: dn}
}

func (s *Sender) SendResetPwdEmail(l mailUsecase.Letter) error {
	tmpl, err := template.ParseFS(letterTmpl, "letter_template.html")
	if err != nil {
		s.log.Err(err).Msg("error parsing letter template from file")
	}
	q := make(url.Values)
	q.Set("token", l.Token)
	// TODO: path should be provided
	u := &url.URL{
		Scheme:   "https",
		Host:     s.dName,
		Path:     "/password/reset",
		RawQuery: q.Encode(),
	}
	link := u.String()
	letter := Letter{
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
	recipients := []mailjet.Recipient{{
		Email: l.To,
		Name:  l.Name,
	}}
	// TODO: refactor this
	email := &mailjet.InfoSendMail{
		FromEmail:  "noreply@traillyst.com",
		FromName:   "Traillyst",
		Subject:    l.Subject,
		TextPart:   l.Header + "\n" + l.Body + "\n" + link,
		HTMLPart:   buf.String(),
		Recipients: recipients,
	}
	_, err = s.mail.SendMail(email)
	if err != nil {
		s.log.Err(err).Msg("error sending email")
		return err
	}
	// TODO: log original trace id
	s.log.Info().Msgf("email sent successfully")
	return nil
}

func (s *Sender) SendRegisterEmail(l mailUsecase.Letter) error {
	tmpl, err := template.ParseFS(letterTmpl, "letter_template.html")
	if err != nil {
		s.log.Err(err).Msg("error parsing letter template from file")
	}
	q := make(url.Values)
	q.Set("token", l.Token)
	q.Set("email", l.To)
	// TODO: path should be provided
	u := &url.URL{
		Scheme:   "https",
		Host:     s.dName,
		Path:     "/user/verify",
		RawQuery: q.Encode(),
	}
	link := u.String()
	letter := Letter{
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
	recipients := []mailjet.Recipient{{
		Email: l.To,
		Name:  l.Name,
	}}
	// TODO: refactor this
	email := &mailjet.InfoSendMail{
		FromEmail:  "noreply@traillyst.com",
		FromName:   "Traillyst",
		Subject:    l.Subject,
		TextPart:   l.Header + "\n" + l.Body + "\n" + link,
		HTMLPart:   buf.String(),
		Recipients: recipients,
	}
	_, err = s.mail.SendMail(email)
	if err != nil {
		s.log.Err(err).Msg("error sending email")
		return err
	}
	// TODO: log original trace id
	s.log.Info().Msgf("email sent successfully")
	return nil
}
