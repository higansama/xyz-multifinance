package events

import (
	"context"
	"sync"
)

type MemoryProvider struct {
	mutex    sync.RWMutex
	Handlers map[string][]Handler
}

func NewMemoryProvider() (*MemoryProvider, error) {
	return &MemoryProvider{}, nil
}

func (mp *MemoryProvider) Publish(ctx context.Context, topic string, e *Event) error {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()

	for _, v := range mp.Handlers[topic] {
		err := v(ctx, *e)

		// since memory provider runs synchronously, we may just immediately
		// return the error if an error occured.
		if err != nil {
			return err
		}
	}

	return nil
}

func (mp *MemoryProvider) Register(opts HandlerOptions, h Handler) {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	if mp.Handlers == nil {
		mp.Handlers = make(map[string][]Handler, 0)
	}

	mp.Handlers[opts.Topic] = append(mp.Handlers[opts.Topic], h)

	return
}

func (mp *MemoryProvider) Subscribe() {
	return
}

func (mp *MemoryProvider) Shutdown() {
	return
}
