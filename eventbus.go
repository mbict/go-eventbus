package eventbus

import (
	"reflect"
)

type EventType string

type Event interface {
	EventType() EventType
}

type EventHandler interface {
	Handle(event Event)
}

type EventHandlerFunc func(event Event)

func (h EventHandlerFunc) Handle(event Event) {
	h(event)
}

type EventBus interface {
	Subscribe(EventHandler, ...EventType)
	Unsubscribe(EventHandler, ...EventType)
	Publish(Event) error
}

func ChainEventHandler(handler EventHandler, wrap EventHandler) EventHandler {
	return EventHandlerFunc(func(event Event) {
		wrap.Handle(event)
		handler.Handle(event)
	})
}

type eventHandlers []EventHandler
type eventChannels map[EventType]eventHandlers

type eventBus struct {
	handlers eventChannels
}

func (eb *eventBus) Subscribe(handler EventHandler, events ...EventType) {
	if len(events) == 0 {
		eb.handlers["*"] = append(eb.handlers["*"], handler)
		return
	}

	for _, eventType := range events {
		eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	}
}

func (eb *eventBus) Unsubscribe(handler EventHandler, events ...EventType) {
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
	for _, handler := range eb.handlers[event.EventType()] {
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
