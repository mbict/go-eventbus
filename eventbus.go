package eventbus

import (
	"reflect"
)

type EventType = string

type Event interface {
	EventType() EventType
}

type EventHandler interface {
	Handle(event Event) error
}

type EventHandlerFunc func(event Event) error

func (h EventHandlerFunc) Handle(event Event) error {
	return h(event)
}

type PublishErrorHandlerFunc func(error, Event) error

type EventBus interface {
	Subscribe(EventHandler, ...EventType)
	Unsubscribe(EventHandler, ...EventType)
	Publish(Event) error
}

func ChainEventHandler(handler EventHandler, wrap EventHandler) EventHandler {
	return EventHandlerFunc(func(event Event) error {
		if err := wrap.Handle(event); err != nil {
			return err
		}
		return handler.Handle(event)
	})
}

type eventHandlers []EventHandler
type eventChannels map[EventType]eventHandlers

type eventBus struct {
	handlers         eventChannels
	errorHandlerFunc PublishErrorHandlerFunc
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
	//specific event name handler
	if err := eb.publishEvent(event, eb.handlers[event.EventType()]); err != nil {
		return err
	}

	//catchall event handlers
	if err := eb.publishEvent(event, eb.handlers["*"]); err != nil {
		return err
	}
	return nil
}

func (eb *eventBus) publishEvent(event Event, handlers eventHandlers) error {
	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			if err := eb.handlePublishError(err, event); err != nil {
				return err
			}
		}
	}
	return nil
}

func (eb *eventBus) handlePublishError(err error, event Event) error {
	if eb.errorHandlerFunc == nil {
		return err
	} else if err := eb.errorHandlerFunc(err, event); err != nil {
		return err
	}
	return nil
}

func New() EventBus {
	return &eventBus{
		handlers: make(eventChannels),
	}
}

func NewWithErrorHandler(errorHandlerFunc PublishErrorHandlerFunc) EventBus {
	return &eventBus{
		handlers:         make(eventChannels),
		errorHandlerFunc: errorHandlerFunc,
	}
}
