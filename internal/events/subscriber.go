package events

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

var wait = make(chan bool)

// Subscribe starts a run loop with a Subscriber that listens to topics and
// waits for a syscall.SIGINT or syscall.SIGTERM
func (eb *EventBus) Subscribe() {
	log.Info().Msgf("Subscribe to event bus: %s", eb.client.ServiceName)
	eb.client.Provider.Subscribe()
	<-wait
}

func (eb *EventBus) Shutdown() {
	log.Info().Msg("Gracefully shutting down subscribers")
	wait <- true
	eb.client.Provider.Shutdown()
}

// HandlerOptions defines the options for a subscriber http
type HandlerOptions struct {
	// The topic to subscribe to
	Topic string
	//	// The name of this subscriber/function
	//	Name string
	// The name of this subscriber/function's service
	ServiceName string
	// The function to invoke
	Handler Handler
	// A message deadline/timeout
	Deadline time.Duration
	// Concurrency sets the maximum number of msgs to be run concurrently
	// default: 20
	Concurrency int
	// Auto Ack the message automatically if return err == nil
	AutoAck bool
	// StartFromBeginning starts a new subscriber from
	// the beginning of messages available, if supported
	StartFromBeginning bool
	// Unique subscriber means that all subscribers will receive all messages
	Unique bool
}

// Listen an event to specific topic
func (eb *EventBus) Listen(opts HandlerOptions) {
	if opts.Topic == "" {
		panic("subscriber@Listen (topic must be set)")
	}

	if opts.Handler == nil {
		panic("subscriber@Listen (http cannot be nil)")
	}

	if opts.ServiceName == "" {
		panic("subscriber@Listen (service name cannot be empty)")
	}

	// Set some default options
	if opts.Deadline == 0 {
		opts.Deadline = 10 * time.Second
	}

	// Set some default concurrency
	if opts.Concurrency == 0 {
		opts.Concurrency = 20
	}

	cb := func(ctx context.Context, m Event) error {
		err := opts.Handler(ctx, m)
		if err != nil {
			err = errors.Errorf("subscriber@Listen (%s): %s", opts.ServiceName, err)
			log.Err(err).Send()
			return err
		}

		return nil
	}

	mw := chainSubscriberMiddleware(eb.client.Middleware...)
	eb.client.Provider.Register(opts, mw(opts, cb))
}

func chainSubscriberMiddleware(mw ...Middleware) func(opts HandlerOptions, next Handler) Handler {
	return func(opts HandlerOptions, final Handler) Handler {
		return func(ctx context.Context, m Event) error {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i].SubscribeInterceptor(opts, last)
			}
			return last(ctx, m)
		}
	}
}
