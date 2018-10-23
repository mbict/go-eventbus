package eventbus

import "testing"

const (
	EventA EventType = "event:test1"
	EventB EventType = "event:test2"
)

type TestEventA struct {
	Handled int
}

func (*TestEventA) EventType() EventType {
	return EventA
}

type TestEventB struct {
	Handled int
}

func (*TestEventB) EventType() EventType {
	return EventB
}

func TestEventBus(t *testing.T) {
	eventHandler := EventHandlerFunc(func(event Event) {
		test := event.(*TestEventA)
		test.Handled++
	})

	bus := New()
	bus.Subscribe(eventHandler, EventA)
	bus.Subscribe(eventHandler, (*TestEventB)(nil).EventType())
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
	eventHandler := EventHandlerFunc(func(event Event) {
		test := event.(*TestEventA)
		test.Handled++
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
	eventHandler := EventHandlerFunc(func(event Event) {
		if v, ok := event.(*TestEventA); ok {
			v.Handled++
		}
		if v, ok := event.(*TestEventB); ok {
			v.Handled++
		}
	})

	bus := New()
	bus.Subscribe(eventHandler, EventA)
	bus.Subscribe(eventHandler, (*TestEventB)(nil).EventType())

	eventA := &TestEventA{}
	eventB := &TestEventB{}

	bus.Unsubscribe(eventHandler, (*TestEventA)(nil).EventType())
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
	eventHandlerA := EventHandlerFunc(func(event Event) {
		if v, ok := event.(*TestEventA); ok {
			v.Handled++
		}
	})
	eventHandlerB := EventHandlerFunc(func(event Event) {
		if v, ok := event.(*TestEventB); ok {
			v.Handled++
		}
	})

	bus := New()
	bus.Subscribe(eventHandlerA, EventA)
	bus.Subscribe(eventHandlerA, (*TestEventB)(nil).EventType())
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
