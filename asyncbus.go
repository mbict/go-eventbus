package eventbus

import (
	"sync"
)

type asyncEventBus struct {
	handlers         map[EventName]eventHandlers
	nameResolver     EventNameResolver
	errorHandlerFunc PublishErrorHandlerFunc
	mu               sync.RWMutex
}

func (eb *asyncEventBus) setEventResolver(resolver EventNameResolver) {
	eb.nameResolver = resolveEventName
}

func (eb *asyncEventBus) setErrorHandler(errorHandler PublishErrorHandlerFunc) {
	eb.errorHandlerFunc = errorHandler
}

func (eb *asyncEventBus) Subscribe(handler EventHandler, events ...EventName) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if len(events) == 0 {
		eb.handlers["*"] = append(eb.handlers["*"], handler)
		return
	}

	for _, eventType := range events {
		eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	}
}

func (eb *asyncEventBus) Unsubscribe(handler EventHandler, events ...EventName) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if len(events) == 0 {
		for eventType := range eb.handlers {
			eb.handlers[eventType] = removeHandlerFromSlice(eb.handlers[eventType], handler)
			if len(eb.handlers[eventType]) == 0 {
				delete(eb.handlers, eventType)
			}
		}
	}

	for _, eventType := range events {
		eb.handlers[eventType] = removeHandlerFromSlice(eb.handlers[eventType], handler)
		if len(eb.handlers[eventType]) == 0 {
			delete(eb.handlers, eventType)
		}
	}
}

func (eb *asyncEventBus) Publish(event any) error {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	for _, handler := range eb.handlers[eb.nameResolver(event)] {
		go func(handler EventHandler) {
			handler.Handle(event)
		}(handler)
	}

	//catchall event handlers
	for _, handler := range eb.handlers["*"] {
		go func(handler EventHandler) {
			handler.Handle(event)
		}(handler)
	}
	return nil
}

func NewAsync(options ...Option) EventBus {
	eb := &asyncEventBus{
		handlers:     make(eventChannels),
		nameResolver: resolveEventName,
	}

	for _, option := range options {
		option(eb)
	}

	return eb
}
