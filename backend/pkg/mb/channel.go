package mb

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Channel struct {
	ch          *amqp.Channel
	QueueName   string
	log         *zerolog.Logger
	isConsuming bool
	sync.Mutex
}

type ChConfig struct {
	QName   string
	WithDLQ bool
}

type Message struct {
	delivery *amqp.Delivery
	Body     []byte
	Type     string
	Headers  amqp.Table
}

func (cm *ConnManager) NewChannel(cfg ChConfig) (*Channel, error) {

	ch, err := cm.conn.Channel()
	if err != nil {
		return nil, err
	}

	args := amqp.Table{}

	if cfg.WithDLQ {
		DLQName := cfg.QName + "_DQL"
		DLXName := cfg.QName + "_DLX"

		// Declare the DLX
		err = ch.ExchangeDeclare(
			DLXName,  // name
			"direct", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			return nil, err
		}

		// Declare the DLQ
		_, err = ch.QueueDeclare(
			DLQName, // name
			true,    // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			return nil, err
		}

		// Bind DLQ to DLX
		err = ch.QueueBind(
			DLQName, // queue name
			"",      // routing key (using the same as the queue name for simplicity)
			DLXName, // exchange
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}

		// Add DLX configuration to arguments for the main queue
		args["x-dead-letter-exchange"] = DLXName
	}

	_, err = ch.QueueDeclare(
		cfg.QName,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return nil, err
	}

	return &Channel{
		ch:          ch,
		QueueName:   cfg.QName,
		log:         cm.log,
		isConsuming: false,
		Mutex:       sync.Mutex{},
	}, nil
}

func (c *Channel) Publish(ctx context.Context, body any) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		c.log.Err(err).Msg("error encoding queue message to JSON")
		return err
	}

	c.Lock()
	defer c.Unlock()

	return c.ch.PublishWithContext(
		ctx,
		"",          // Exchange
		c.QueueName, // Routing key
		false,       // Mandatory
		false,       // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		})
}

func (c *Channel) Consume() (<-chan Message, error) {

	if c.isConsuming {
		e := errors.New("consumer has already been registered for this channel")
		c.log.Err(e).Msg("error register consumer for queue")
		return nil, e
	}

	msgs, err := c.ch.Consume(
		c.QueueName,
		"",    // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)
	if err != nil {
		c.log.Err(err).Msg("error initializing read channel for queue")
		return nil, err
	}

	out := make(chan Message)

	go func() {
		defer close(out)

		for d := range msgs {
			data := Message{
				delivery: &d,
				Body:     d.Body,
				Type:     "application/json",
				Headers:  d.Headers,
			}
			out <- data
		}
	}()
	return out, nil
}

func (m *Message) Ack(multiple bool) error {
	return m.delivery.Ack(multiple)
}

func (m *Message) Nack(multiple bool, requeue bool) error {
	return m.delivery.Nack(multiple, requeue)
}

func (m *Message) Reject(requeue bool) error {
	return m.delivery.Reject(requeue)
}

func (c *Channel) Close() {
	if err := c.ch.Close(); err != nil {
		c.log.Err(err).Msg("failed to close queue channel")
	}
}
