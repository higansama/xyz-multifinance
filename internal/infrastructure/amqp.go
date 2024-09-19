package infrastructure

import (
	"fmt"
	"net"

	"github.com/higansama/xyz-multi-finance/config"
	"github.com/pkg/errors"
	amqpgo "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func NewAmqp(cfg config.RabbitMQ) (*amqpgo.Connection, error) {
	uri := cfg.URI

	if uri == "" {
		scheme := "amqp"
		if cfg.Secure {
			scheme = "amqps"
		}

		uri = fmt.Sprintf(
			"%s://%s:%s@%s",
			scheme,
			cfg.Username,
			cfg.Password,
			net.JoinHostPort(cfg.Host, cfg.Port))
	}

	conn, err := amqpgo.Dial(uri)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return conn, nil
}

func (infra *Infrastructure) configAmqp() (*amqpgo.Connection, error, func()) {
	conn, err := NewAmqp(infra.Config.Messaging.RabbitMQ)
	if err != nil {
		return nil, err, nil
	}

	log.Info().Msg("Amqp connected")

	return conn, nil, func() {
		_ = conn.Close()
	}
}
