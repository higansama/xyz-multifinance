package amqp

import (
	"context"
	"fmt"
	"math"
	"time"

	ierrors "github.com/higansama/xyz-multi-finance/internal/errors"
	"github.com/higansama/xyz-multi-finance/internal/inbox"
	"github.com/higansama/xyz-multi-finance/internal/infrastructure"
	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/pkg/errors"
	amqpgo "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type ConsumerConfiguration struct {
	Exchange       string
	Queue          amqpgo.Queue
	Dlx            string
	Dlq            amqpgo.Queue
	Plq            amqpgo.Queue
	StrictOrdering bool
}

func SetupConsumer(
	infra *infrastructure.Infrastructure,
	exchange string,
	topic string,
	consumerIdentifier string,
	strictOrdering bool,
) (*ConsumerConfiguration, error) {
	exchange = utils.FormatNameForEnv(infra.Env, exchange)

	channel, err := infra.Amqp.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed opening AMQP channel")
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		exchange,
		amqpgo.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed declaring exchange")
	}

	queueName := utils.FormatQueueName(
		infra.Env,
		infra.Config.App.Name,
		fmt.Sprintf("%s.%s", consumerIdentifier, topic))
	if topic == "#" || topic == "*" {
		queueName = utils.FormatQueueName(
			infra.Env,
			infra.Config.App.Name,
			fmt.Sprintf("%s.%s.%s", consumerIdentifier, exchange, topic))
	}

	dlx := exchange + ".dlx"      // dead letter exchange
	dlqName := queueName + ".dlq" // dead letter queue
	plqName := queueName + ".plq" // parking lot queue

	qheaders := amqpgo.Table{}
	if strictOrdering {
		// simple way to maintain message processing order
		qheaders["x-single-active-consumer"] = true
	} else {
		err = channel.ExchangeDeclare(
			dlx,
			amqpgo.ExchangeTopic,
			true,
			false,
			false,
			false,
			nil)
		if err != nil {
			return nil, fmt.Errorf("failed declaring dlx: %w", err)
		}

		qheaders["x-dead-letter-exchange"] = dlx
	}

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		qheaders)
	if err != nil {
		return nil, errors.Wrap(err, "failed declaring queue")
	}

	err = channel.QueueBind(
		queue.Name,
		topic,
		exchange,
		false,
		nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed bind queue to exchange")
	}

	cconf := &ConsumerConfiguration{
		Exchange:       exchange,
		Queue:          queue,
		StrictOrdering: strictOrdering,
	}

	if !strictOrdering {
		dlq, err := channel.QueueDeclare(
			dlqName,
			true,
			false,
			false,
			false,
			amqpgo.Table{
				//"x-single-active-consumer": true,
				"x-dead-letter-exchange": exchange,
				// message in dlq will "dead lettered" to original exchange
				"x-message-ttl": (1 * time.Second).Milliseconds(),
				// how long to push back message to original queue
			})
		if err != nil {
			return nil, fmt.Errorf("failed declaring dlq: %w", err)
		}

		err = channel.QueueBind(
			dlq.Name,
			topic,
			dlx,
			false,
			nil)
		if err != nil {
			return nil, fmt.Errorf("failed bind dlq to dlx: %w", err)
		}

		plq, err := channel.QueueDeclare(
			plqName,
			true,
			false,
			false,
			false,
			amqpgo.Table{
				//"x-single-active-consumer": true,
			})
		if err != nil {
			return nil, fmt.Errorf("failed declaring plq: %w", err)
		}

		cconf.Dlx = dlx
		cconf.Dlq = dlq
		cconf.Plq = plq
	}

	return cconf, nil
}

func StartConsumer(
	infra *infrastructure.Infrastructure,
	cconf *ConsumerConfiguration,
	errCh chan error,
	consumer string,
	handler func(ctx context.Context, delivery amqpgo.Delivery) error,
) {
	amqpChannel, err := infra.Amqp.Channel()
	if err != nil {
		errCh <- errors.Wrapf(err, "[%s] failed opening AMQP channel", cconf.Queue.Name)
		return
	}
	defer amqpChannel.Close()

	deliveries, err := amqpChannel.Consume(
		cconf.Queue.Name,
		consumer,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		errCh <- errors.Wrapf(err, "[%s] failed starting AMQP consumer", cconf.Queue.Name)
		return
	}

	chClosed := make(chan *amqpgo.Error)
	amqpChannel.NotifyClose(chClosed)

	for {
		select {
		case <-infra.Ctx.Done():
			return
		case _ = <-chClosed:
			errCh <- errors.New("channel closed")
			return
		case d := <-deliveries:
			fn := func() {
				xDeath := getXDeathCount(d)
				maxRetry := int64(3)

				if xDeath > 0 {
					time.Sleep(time.Duration(int64(math.Pow(3, float64(xDeath-1))) * (10 * time.Second).Nanoseconds())) // exponential backoff
					log.Info().Msgf("%dx retry processing message %s", xDeath, d.MessageId)
				} else {
					log.Info().Msgf("Processing message %s", d.MessageId)
				}

				inboxItem, err := infra.InboxRepo.
					FindByMessageIdAndConsumer(infra.Ctx, d.MessageId, consumer)
				if err != nil {
					errCh <- err
					return
				}
				if inboxItem != nil { // idempotency
					err = d.Ack(false)
					if err != nil {
						errCh <- errors.WithMessage(err, "error ack for idempotency")
						return
					}
					return
				}

				err = infra.TxManager.WithTransaction(infra.Ctx, func(txCtx context.Context) error {
					// indicates that message have been proceeded.
					err = infra.InboxRepo.Save(txCtx, inbox.New(d.MessageId, consumer))
					if err != nil {
						return err
					}

					err = handler(txCtx, d)
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					if cconf.StrictOrdering {
						errCh <- errors.WithMessage(err, "error when executing http")
						return
					} else {
						if xDeath < maxRetry {
							if xDeath == 0 {
								log.Err(err).Send()
								log.Info().Msgf("... Dead lettered message %s", d.MessageId)
							}
							// send to dlx
							err = d.Nack(false, false)
							if err != nil {
								errCh <- ierrors.RecoveredError{ActualErr: errors.WithStack(err)}
							}

							return // prevent ack
						} else {
							log.Err(err).Send()
							outerErr := err
							// send to plq
							log.Info().Msgf("Send message %s to plq", d.MessageId)
							sendPlq := func() error {
								channel, err := infra.Amqp.Channel()
								if err != nil {
									return errors.Wrap(err, "failed opening AMQP channel")
								}
								defer channel.Close()

								err = channel.PublishWithContext(
									infra.Ctx,
									"",
									cconf.Plq.Name,
									false,
									false,
									deliveryToPublishing(d, outerErr))
								if err != nil {
									return errors.Wrap(err, "failed publish to PLQ")
								}

								return nil
							}

							err = sendPlq()
							if err != nil {
								errCh <- ierrors.RecoveredError{ActualErr: errors.WithStack(err)}
								return
							}
						}
					}
				}

				err = d.Ack(false)
				if err != nil {
					errCh <- ierrors.RecoveredError{
						ActualErr: errors.WithMessage(err, "error when ack the delivery"),
					}
					return
				}
			}

			if cconf.StrictOrdering {
				fn()
			} else {
				go fn()
			}
		}
	}
}

func getXDeathCount(delivery amqpgo.Delivery) int64 {
	if v, ok := delivery.Headers["x-death"]; ok {
		if v, ok := v.([]interface{}); ok {
			if len(v) > 0 {
				t := v[0].(amqpgo.Table)
				if v, ok := t["count"]; ok {
					return v.(int64)
				}
			}
		}
	}

	return 0
}

func deliveryToPublishing(delivery amqpgo.Delivery, err error) amqpgo.Publishing {
	var headers amqpgo.Table

	if err != nil {
		headers = modifyAmqpTable(headers, func(h amqpgo.Table) {
			h["error"] = err.Error()
		})
	}

	return amqpgo.Publishing{
		Headers:      headers,
		ContentType:  delivery.ContentType,
		Timestamp:    delivery.Timestamp,
		MessageId:    delivery.MessageId,
		Body:         delivery.Body,
		DeliveryMode: delivery.DeliveryMode,
		Priority:     delivery.Priority,
		AppId:        delivery.AppId,
	}
}

func modifyAmqpTable(table amqpgo.Table, fn func(table amqpgo.Table)) amqpgo.Table {
	if table == nil {
		table = amqpgo.Table{}
	}
	fn(table)
	return table
}
