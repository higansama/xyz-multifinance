package events

import (
	"context"
)

// NoopProvider is a simple provider that does nothing
type NoopProvider struct{}

// Publish does nothing
func (np NoopProvider) Publish(ctx context.Context, topic string, e *Event) error {
	return nil
}

// Register does nothing
func (np NoopProvider) Register(opts HandlerOptions, h Handler) {
	return
}

// Subscribe does nothing
func (np NoopProvider) Subscribe(opts HandlerOptions, h Handler) {
	return
}

// Shutdown shutsdown immediately
func (np NoopProvider) Shutdown() {
	return
}
