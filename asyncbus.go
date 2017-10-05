package eventbus

import (
	"sync"
)

type asyncEventBus struct {
	handlers map[string]eventHandlers
	mu       sync.RWMutex
}

func (eb *asyncEventBus) Subscribe(handler EventHandler, events ...Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if len(events) == 0 {
		eb.handlers["*"] = append(eb.handlers["*"], handler)
		return
	}

	for _, event := range events {
		eb.handlers[event.EventName()] = append(eb.handlers[event.EventName()], handler)
	}
}

func (eb *asyncEventBus) Unsubscribe(handler EventHandler, events ...Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if len(events) == 0 {
		for k := range eb.handlers {
			eb.handlers[k] = removeHandlerFromSlice(eb.handlers[k], handler)
			if len(eb.handlers[k]) == 0 {
				delete(eb.handlers, k)
			}
		}
	}

	for _, e := range events {
		k := e.EventName()
		eb.handlers[k] = removeHandlerFromSlice(eb.handlers[k], handler)
		if len(eb.handlers[k]) == 0 {
			delete(eb.handlers, k)
		}
	}
}

func (eb *asyncEventBus) Publish(event Event) error {
	eb.mu.RLock()
	var wg sync.WaitGroup

	for _, handler := range eb.handlers[event.EventName()] {
		wg.Add(1)
		go func(handler EventHandler) {
			handler.Handle(event)
			wg.Done()
		}(handler)
	}

	//catchall event handlers
	for _, handler := range eb.handlers["*"] {
		wg.Add(1)
		go func(handler EventHandler) {
			handler.Handle(event)
			wg.Done()
		}(handler)
	}

	eb.mu.RUnlock()
	wg.Wait()
	return nil
}

func NewAsync() EventBus {
	return &asyncEventBus{
		handlers: make(eventChannels),
	}
}
