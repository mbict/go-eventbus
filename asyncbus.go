package eventbus

import (
	"sync"
)

type asyncEventBus struct {
	handlers         map[EventType]eventHandlers
	errorHandlerFunc PublishErrorHandlerFunc
	mu               sync.Mutex
}

func (eb *asyncEventBus) Subscribe(handler EventHandler, events ...EventType) {
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

func (eb *asyncEventBus) Unsubscribe(handler EventHandler, events ...EventType) {
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

func (eb *asyncEventBus) Publish(event Event) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	var wg sync.WaitGroup

	for _, handler := range eb.handlers[event.EventType()] {
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

	wg.Wait()
	return nil
}

func NewAsync() EventBus {
	return &asyncEventBus{
		handlers: make(eventChannels),
	}
}

func NewAsyncWithErrorHandler(errorHandlerFunc PublishErrorHandlerFunc) EventBus {
	return &asyncEventBus{
		handlers:         make(eventChannels),
		errorHandlerFunc: errorHandlerFunc,
	}
}
