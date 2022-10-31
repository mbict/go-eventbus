package eventbus

import (
	"errors"
	"testing"
)

const (
	EventA EventName = "event:test1"
	EventB EventName = "event:test2"
)

type TestEventA struct {
	Handled int
}

func (*TestEventA) EventName() EventName {
	return EventA
}

type TestEventB struct {
	Handled int
}

func (*TestEventB) EventName() EventName {
	return EventB
}

func TestEventBus(t *testing.T) {
	eventHandler := EventHandlerFunc(func(event any) error {
		test := event.(*TestEventA)
		test.Handled++
		return nil
	})

	bus := New()
	bus.Subscribe(eventHandler, EventA)
	bus.Subscribe(eventHandler, (*TestEventB)(nil).EventName())
	event := &TestEventA{}

	err := bus.Publish(event)

	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	if event.Handled != 1 {
		t.Fatalf("expected event to be handled once, but is handled %d times", event.Handled)
	}
}

func TestEventBusCatchAll(t *testing.T) {
	eventHandler := EventHandlerFunc(func(event any) error {
		test := event.(*TestEventA)
		test.Handled++
		return nil
	})

	bus := New()
	bus.Subscribe(eventHandler)
	bus.Subscribe(eventHandler, EventB)
	event := &TestEventA{}

	err := bus.Publish(event)

	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	if event.Handled != 1 {
		t.Fatalf("expected event to be handled once, but is handled %d times", event.Handled)
	}
}

func TestEventBusUnsubscribeForSpecficEvent(t *testing.T) {
	eventHandler := EventHandlerFunc(func(event any) error {
		if v, ok := event.(*TestEventA); ok {
			v.Handled++
		}
		if v, ok := event.(*TestEventB); ok {
			v.Handled++
		}
		return nil
	})

	bus := New()
	bus.Subscribe(eventHandler, EventA)
	bus.Subscribe(eventHandler, (*TestEventB)(nil).EventName())

	eventA := &TestEventA{}
	eventB := &TestEventB{}

	bus.Unsubscribe(eventHandler, (*TestEventA)(nil).EventName())
	err := bus.Publish(eventA)

	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	if eventA.Handled != 0 {
		t.Fatalf("expected event A to be unhandled, but is handled %d times", eventA.Handled)
	}

	err = bus.Publish(eventB)
	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	if eventB.Handled != 1 {
		t.Fatalf("expected event to be handled once, but is handled %d times", eventB.Handled)
	}
}

func TestEventBusUnsubscribeAll(t *testing.T) {
	eventHandlerA := EventHandlerFunc(func(event any) error {
		if v, ok := event.(*TestEventA); ok {
			v.Handled++
		}
		return nil
	})
	eventHandlerB := EventHandlerFunc(func(event any) error {
		if v, ok := event.(*TestEventB); ok {
			v.Handled++
		}
		return nil
	})

	bus := New()
	bus.Subscribe(eventHandlerA, EventA)
	bus.Subscribe(eventHandlerA, (*TestEventB)(nil).EventName())
	bus.Subscribe(eventHandlerB, EventB)
	bus.Subscribe(eventHandlerA)
	bus.Subscribe(eventHandlerB)

	eventA := &TestEventA{}
	eventB := &TestEventB{}

	bus.Unsubscribe(eventHandlerA)
	bus.Publish(eventA)
	bus.Publish(eventB)

	if eventA.Handled != 0 {
		t.Fatalf("expected event A to be unhandled, but is handled %d times", eventA.Handled)
	}

	if eventB.Handled != 2 {
		t.Fatalf("expected event to be handled once, but is handled %d times", eventB.Handled)
	}
}

func TestEventBusPublishErrorWithoutErrorHandler(t *testing.T) {
	expectedErr := errors.New("unhandled error")
	eventHandler := EventHandlerFunc(func(event any) error {
		return expectedErr
	})

	bus := New()
	bus.Subscribe(eventHandler, EventA)
	event := &TestEventA{}

	err := bus.Publish(event)

	if err != expectedErr {
		t.Fatalf("expected error, but got %v", err)
	}
}

func TestEventBusHandlesError(t *testing.T) {
	expectedErr := errors.New("unhandled error")
	expectedEvent := &TestEventA{}
	eventHandler := EventHandlerFunc(func(event any) error {
		if v, ok := event.(*TestEventA); ok {
			v.Handled++
		}
		return expectedErr
	})

	bus := New(WithErrorHandler(func(err error, event any) error {
		if err != expectedErr {
			t.Fatalf("expected error, but got %v", err)
		}
		return nil
	}))
	bus.Subscribe(eventHandler, EventA)

	err := bus.Publish(expectedEvent)

	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	if expectedEvent.Handled != 1 {
		t.Fatalf("expected event to be handled once, but is handled %d times", expectedEvent.Handled)
	}
}

func TestEventBusReturnsErrorWhenPublishErrorHandlerFails(t *testing.T) {
	expectedErr := errors.New("unhandled error")
	expectedPublishErr := errors.New("unhandled publish error")
	expectedEvent := &TestEventA{}
	eventHandler := EventHandlerFunc(func(event any) error {
		if v, ok := event.(*TestEventA); ok {
			v.Handled++
		}
		return expectedErr
	})

	bus := New(WithErrorHandler(func(err error, event any) error {
		if err != expectedErr {
			t.Fatalf("expected error, but got %v", err)
		}
		return expectedPublishErr
	}))
	bus.Subscribe(eventHandler, EventA)

	err := bus.Publish(expectedEvent)

	if err != expectedPublishErr {
		t.Fatalf("expected error, but got %v", expectedPublishErr)
	}

	if expectedEvent.Handled != 1 {
		t.Fatalf("expected event to be handled once, but is handled %d times", expectedEvent.Handled)
	}
}
