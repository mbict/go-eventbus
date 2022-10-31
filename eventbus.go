package eventbus

import (
	"reflect"
)

type EventName = string

type Event interface {
	EventName() EventName
}

type EventHandler interface {
	Handle(event any) error
}

type EventHandlerFunc func(event any) error

func (h EventHandlerFunc) Handle(event any) error {
	return h(event)
}

type PublishErrorHandlerFunc func(error, any) error

type EventBus interface {
	Subscribe(EventHandler, ...EventName)
	Unsubscribe(EventHandler, ...EventName)
	Publish(any) error
}

func Chain(handler EventHandler, wrap EventHandler) EventHandler {
	return EventHandlerFunc(func(event any) error {
		if err := wrap.Handle(event); err != nil {
			return err
		}
		return handler.Handle(event)
	})
}

type eventHandlers []EventHandler
type eventChannels map[EventName]eventHandlers

type eventBus struct {
	handlers          eventChannels
	errorHandlerFunc  PublishErrorHandlerFunc
	eventNameResolver EventNameResolver
}

func (eb *eventBus) setEventResolver(resolver EventNameResolver) {
	eb.eventNameResolver = resolver
}

func (eb *eventBus) setErrorHandler(errorHandler PublishErrorHandlerFunc) {
	eb.errorHandlerFunc = errorHandler
}

func (eb *eventBus) Subscribe(handler EventHandler, events ...EventName) {
	if len(events) == 0 {
		eb.handlers["*"] = append(eb.handlers["*"], handler)
		return
	}

	for _, eventType := range events {
		eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	}
}

func (eb *eventBus) Unsubscribe(handler EventHandler, events ...EventName) {
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

func (eb *eventBus) Publish(event any) error {
	//specific event name handler
	if err := eb.publishEvent(event, eb.handlers[eb.eventNameResolver(event)]); err != nil {
		return err
	}

	//catchall event handlers
	if err := eb.publishEvent(event, eb.handlers["*"]); err != nil {
		return err
	}
	return nil
}

func (eb *eventBus) publishEvent(event any, handlers eventHandlers) error {
	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			if err := eb.handlePublishError(err, event); err != nil {
				return err
			}
		}
	}
	return nil
}

func (eb *eventBus) handlePublishError(err error, event any) error {
	if eb.errorHandlerFunc == nil {
		return err
	} else if err := eb.errorHandlerFunc(err, event); err != nil {
		return err
	}
	return nil
}

func New(options ...Option) EventBus {
	eb := &eventBus{
		handlers:          make(eventChannels),
		eventNameResolver: resolveEventName,
	}

	for _, option := range options {
		option(eb)
	}

	return eb
}
