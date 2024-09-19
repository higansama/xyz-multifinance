package events

import (
	"context"
	"time"
)

// Handler is the message http
type Handler func(ctx context.Context, e Event) error

// PublishHandler wraps a call to publish, for interception
type PublishHandler func(ctx context.Context, topic string, m *Event) error

// Middleware is an interface to provide subscriber and publisher interceptors
type Middleware interface {
	SubscribeInterceptor(opts HandlerOptions, next Handler) Handler
	PublisherMsgInterceptor(serviceName string, next PublishHandler) PublishHandler
}

// Event is representation of the messaging entry
type Event struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Source      string            `json:"source"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Data        []byte            `json:"data"`
	PublishedAt time.Time         `json:"published_at"`
	Ack         func()            `json:"-"`
	Nack        func()            `json:"-"`
}

// Provider is generic interface for a event bus provider
type Provider interface {
	Publish(ctx context.Context, topic string, m *Event) error
	Register(opts HandlerOptions, handler Handler)
	Subscribe()
	Shutdown()
}

// Client holds a reference to a Provider
type Client struct {
	ServiceName string
	Provider    Provider
	Middleware  []Middleware
}

type EventBus struct {
	client *Client
}

// NewEventBus creates new event bus
func NewEventBus(c *Client) *EventBus {
	return &EventBus{
		client: c,
	}
}
