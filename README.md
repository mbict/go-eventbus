[![Build Status](https://travis-ci.org/mbict/go-eventbus.png?branch=master)](https://travis-ci.org/mbict/go-eventbus)
[![GoDoc](https://godoc.org/github.com/mbict/go-eventbus?status.png)](http://godoc.org/github.com/mbict/go-eventbus)
[![GoCover](http://gocover.io/_badge/github.com/mbict/go-eventbus)](http://gocover.io/github.com/mbict/go-eventbus)
[![GoReportCard](http://goreportcard.com/badge/mbict/go-eventbus)](http://goreportcard.com/report/mbict/go-eventbus)

# Event Bus

A simple implementation of a eventbus for go.
Also added a async and concurrent version of the eventbus.

#### Example
A complete example
```go
package main

import (
	"fmt"
	eb "github.com/mbict/go-eventbus"
)

// MyEvent
type MyEvent struct {
	Message string
}

// EventName is needed to identify this event, Pointer receiver
// Reflection is just too expensive for this
func (*MyEvent) EventName() string {
	return "my.event"
}

type OtherEvent struct {
	Message string
}

// EventName is needed to identify this event
// Reflection is just too expensive for this
func (OtherEvent) EventName() string {
	return "other.event"
}

func main() {
	//example of the event handler with a pointer receiver
	eventHandler := eb.EventHandlerFunc(func(event eb.Event) {
		e := event.(*MyEvent)
		fmt.Println("handled event", e.Message)
	})

	//example of the event handler
	otherEventHandler := eb.EventHandlerFunc(func(event eb.Event) {
		e := event.(OtherEvent)
		fmt.Println("handled event", e.Message)
	})

	//wildcard handler
	catchallEventHandler := eb.EventHandlerFunc(func(event eb.Event) {
		switch e := event.(type) {
		case *MyEvent:
			fmt.Println("my event triggered in catch all", e.Message)
		case OtherEvent:
			fmt.Println("other event triggered in catch all", e.Message)
		}
	})

	bus := eb.New()

	//subscribe with a pointer event
	bus.Subscribe(eventHandler, (*MyEvent)(nil))

	//subscribe with a event
	bus.Subscribe(otherEventHandler, OtherEvent{})

	//subscribe to all events
	bus.Subscribe(catchallEventHandler)

	// create a event and send it to the bus
	event1 := &MyEvent{
		Message: "hello you!",
	}
	bus.Publish(event1)

	// and now the normal receiver
	event2 := OtherEvent{
		Message: "hey you are the other one, also hello to you too!",
	}
	bus.Publish(event2)

	//you can also unsubscribe the handler for specific events (if registered)
	bus.Unsubscribe(eventHandler, (*MyEvent)(nil))
	bus.Unsubscribe(otherEventHandler, OtherEvent{})

	//or from all the events, also specific events will be unsubscribed if the handler matches
	bus.Unsubscribe(eventHandler)
}
```

#### Todo
- Test concurrent eventbus
- Test async eventbus