package queue

import (
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type ConnManager struct {
	log  *zerolog.Logger
	conn *amqp.Connection
}

type ConnConfig struct {
	User     string
	Password string
	Host     string
	Log      *zerolog.Logger
}

func NewManager(cfg ConnConfig) (*ConnManager, error) {
	u := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   cfg.Host,
	}

	cs := u.String()

	conn, err := amqp.Dial(cs)
	if err != nil {
		cfg.Log.Err(err).Msg("failed to connect to rabbitmq")
		return nil, err
	}

	return &ConnManager{
		log:  cfg.Log,
		conn: conn,
	}, nil

}

func (cm ConnManager) Close() {
	if err := cm.conn.Close(); err != nil {
		cm.log.Err(err).Msg("failed to close queue connection")
	}
}
