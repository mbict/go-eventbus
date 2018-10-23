package eventbus

import (
	"sync"
)

type concurrentEventBus struct {
	EventBus
	mu sync.Mutex
}

func (eb *concurrentEventBus) Subscribe(handler EventHandler, events ...EventType) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.EventBus.Subscribe(handler, events...)
}

func (eb *concurrentEventBus) Unsubscribe(handler EventHandler, events ...EventType) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.EventBus.Unsubscribe(handler, events...)
}

func (eb *concurrentEventBus) Publish(event Event) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	return eb.EventBus.Publish(event)
}

func NewConcurrent() EventBus {
	return &concurrentEventBus{
		EventBus: New(),
	}
}
