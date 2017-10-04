package eventbus

import (
	"reflect"
	"sync"
)

type Event interface {
	EventName() string
}

type EventHandler interface {
	Handle(event Event)
}

type EventHandlerFunc func(event Event)

func (h EventHandlerFunc) Handle(event Event) {
	h(event)
}

type EventBus interface {
	Subscribe(EventHandler, ...Event)
	Unsubscribe(EventHandler, ...Event)
	Publish(Event) error
}

func ChainEventHandler(handler EventHandler, wrap EventHandler) EventHandler {
	return EventHandlerFunc(func(event Event) {
		wrap.Handle(event)
		handler.Handle(event)
	})
}

type eventHandlers []EventHandler
type eventChannels map[string]eventHandlers

type eventBus struct {
	handlers eventChannels
}

func (eb *eventBus) Subscribe(handler EventHandler, events ...Event) {
	if len(events) == 0 {
		eb.handlers["*"] = append(eb.handlers["*"], handler)
		return
	}

	for _, event := range events {
		eb.handlers[event.EventName()] = append(eb.handlers[event.EventName()], handler)
	}
}

func (eb *eventBus) Unsubscribe(handler EventHandler, events ...Event) {
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

func removeHandlerFromSlice(eh eventHandlers, h EventHandler) eventHandlers {
	for i := len(eh) - 1; i >= 0; i-- {
		if reflect.ValueOf(eh[i]) == reflect.ValueOf(h) {
			copy(eh[i:], eh[i+1:])
			eh[len(eh)-1] = nil
			eh = eh[:len(eh)-1]
		}
	}
	return eh
}

func (eb *eventBus) Publish(event Event) error {
	for _, handler := range eb.handlers[event.EventName()] {
		handler.Handle(event)
	}

	//catchall event handlers
	for _, handler := range eb.handlers["*"] {
		handler.Handle(event)
	}
	return nil
}

func New() EventBus {
	return &eventBus{
		handlers: make(eventChannels),
	}
}

type concurrentEventBus struct {
	EventBus
	mu sync.RWMutex
}

func (eb *concurrentEventBus) Subscribe(handler EventHandler, events ...Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.EventBus.Subscribe(handler, events...)
}

func (eb *concurrentEventBus) Unsubscribe(handler EventHandler, events ...Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.EventBus.Unsubscribe(handler, events...)
}

func (eb *concurrentEventBus) Publish(event Event) error {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return eb.EventBus.Publish(event)
}

func NewConcurrent() EventBus {
	return &concurrentEventBus{
		EventBus: New(),
	}
}

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
