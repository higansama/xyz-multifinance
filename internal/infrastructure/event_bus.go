package infrastructure

import (
	"github.com/higansama/xyz-multi-finance/internal/events"
)

func (infra *Infrastructure) setupEventBus() (*events.EventBus, error) {
	provider, err := events.NewMemoryProvider()
	if err != nil {
		return nil, err
	}

	eb := events.NewEventBus(&events.Client{
		ServiceName: "memory",
		Provider:    provider,
		Middleware:  []events.Middleware{},
	})

	return eb, nil
}
