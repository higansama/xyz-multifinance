package events

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

// Publish published on the client
func (c *Client) Publish(ctx context.Context, topic string, msg any) error {
	b, err := json.Marshal(msg)
	eventType := ""
	t := reflect.TypeOf(msg)
	if t.Kind() == reflect.Ptr {
		eventType = t.Elem().Name()
	} else if t.Kind() == reflect.Struct {
		eventType = t.Name()
	}

	if err != nil {
		return errors.WithStack(err)
	}

	m := &Event{
		ID:          uuid.New().String(),
		Type:        eventType,
		Data:        b,
		PublishedAt: time.Now(),
	}

	mw := chainPublisherMiddleware(c.Middleware...)
	return mw(c.ServiceName, func(ctx context.Context, topic string, m *Event) error {
		return c.Provider.Publish(ctx, topic, m)
	})(ctx, topic, m)
}

func (eb *EventBus) Publish(ctx context.Context, topic string, obj any) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)
		err := eb.client.Publish(ctx, topic, obj)
		if err != nil {
			err = errors.WithStack(err)
		}

		errCh <- err
	}()

	return errCh
}

func chainPublisherMiddleware(mw ...Middleware) func(serviceName string, next PublishHandler) PublishHandler {
	return func(serviceName string, final PublishHandler) PublishHandler {
		return func(ctx context.Context, topic string, m *Event) error {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i].PublisherMsgInterceptor(serviceName, last)
			}
			return last(ctx, topic, m)
		}
	}
}
