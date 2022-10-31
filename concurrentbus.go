package eventbus

import (
	"sync"
)

type concurrentEventBus struct {
	EventBus
	mu sync.Mutex
}

func (eb *concurrentEventBus) Subscribe(handler EventHandler, events ...EventName) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.EventBus.Subscribe(handler, events...)
}

func (eb *concurrentEventBus) Unsubscribe(handler EventHandler, events ...EventName) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.EventBus.Unsubscribe(handler, events...)
}

func (eb *concurrentEventBus) Publish(event any) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	return eb.EventBus.Publish(event)
}

func NewConcurrent(options ...Option) EventBus {
	return &concurrentEventBus{
		EventBus: New(options...),
	}
}
