package main

import (
	"fmt"
	eb "github.com/mbict/go-eventbus/v2"
)

// MyEvent is the event type descriptor
const MyEvent eb.EventName = "my.event"

// MyEventPayload holds the event data
type MyEventPayload struct {
	Message string
}

// EventName is needed to identify this event
func (*MyEventPayload) EventName() eb.EventName {
	return MyEvent
}

// OtherEvent is the event type desciptor
const OtherEvent eb.EventName = "other.event"

// OtherEventPayload holds the event data
type OtherEventPayload struct {
	Message string
}

// EventName is needed to identify this event
func (OtherEventPayload) EventName() eb.EventName {
	return OtherEvent
}

func main() {
	//example of the event handler with a pointer receiver
	eventHandler := eb.EventHandlerFunc(func(event any) error {
		e := event.(*MyEventPayload)
		fmt.Println("handled event", e.Message)
		return nil
	})

	//example of the event handler
	otherEventHandler := eb.EventHandlerFunc(func(event any) error {
		e := event.(OtherEventPayload)
		fmt.Println("handled event", e.Message)
		return nil
	})

	//wildcard handler
	catchallEventHandler := eb.EventHandlerFunc(func(event any) error {
		switch e := event.(type) {
		case *MyEventPayload:
			fmt.Println("my event triggered in catch all", e.Message)
		case OtherEventPayload:
			fmt.Println("other event triggered in catch all", e.Message)
		}
		return nil
	})

	bus := eb.New()

	//subscribe with a pointer event
	bus.Subscribe(eventHandler, MyEvent)

	//subscribe with a event
	bus.Subscribe(otherEventHandler, OtherEvent)

	//subscribe to all events
	bus.Subscribe(catchallEventHandler)

	// create a event and send it to the bus
	event1 := &MyEventPayload{
		Message: "hello you!",
	}
	bus.Publish(event1)

	// and now the normal receiver
	event2 := OtherEventPayload{
		Message: "hey you are the other one, also hello to you too!",
	}
	bus.Publish(event2)

	//you can also unsubscribe the handler for specific events (if registered)
	bus.Unsubscribe(eventHandler, MyEvent)
	bus.Unsubscribe(otherEventHandler, OtherEvent)

	//or from all the events, also specific events will be unsubscribed if the handler matches
	bus.Unsubscribe(eventHandler)
}
