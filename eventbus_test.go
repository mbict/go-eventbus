package eventbus

import "testing"

type TestEventA struct {
	Handled int
}

func (*TestEventA) EventName() string {
	return "test.event1"
}

type TestEventB struct {
	Handled int
}

func (*TestEventB) EventName() string {
	return "test.event2"
}

func TestEventBus(t *testing.T) {
	eventHandler := EventHandlerFunc(func(event Event) {
		test := event.(*TestEventA)
		test.Handled++
	})

	bus := New()
	bus.Subscribe(eventHandler, (*TestEventA)(nil))
	bus.Subscribe(eventHandler, (*TestEventB)(nil))
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
	bus.Subscribe(eventHandler, (*TestEventB)(nil))
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
	bus.Subscribe(eventHandler, (*TestEventA)(nil))
	bus.Subscribe(eventHandler, (*TestEventB)(nil))

	eventA := &TestEventA{}
	eventB := &TestEventB{}

	bus.Unsubscribe(eventHandler, (*TestEventA)(nil))
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
	bus.Subscribe(eventHandlerA, (*TestEventA)(nil))
	bus.Subscribe(eventHandlerA, (*TestEventB)(nil))
	bus.Subscribe(eventHandlerB, (*TestEventB)(nil))
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
