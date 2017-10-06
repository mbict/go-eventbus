package main

import (
	"fmt"
	"time"
	eventbus "github.com/mbict/go-eventbus"
)

type testEvent struct {
	Number int
}

func (testEvent) EventName() string {
	return "testEvent"
}

var counter int

func handler(event eventbus.Event) {
	counter++
}

//just a simple brute force test
func main() {
	eb, cancelFunc := eventbus.NewChanneldWith( eventbus.New() )
	eb.Subscribe( eventbus.EventHandlerFunc(handler) )

	eb.Publish(testEvent{Number:9999999999})
	go func() {
		for i := 0; i < 100000; i++ {
				eb.Publish(testEvent{Number: i})
		}
	}()
	go func() {
		for i := 100000; i < 200000; i++ {
			eb.Publish(testEvent{Number: i})
		}
	}()
	go func() {
		for i := 200000; i < 300000; i++ {
			eb.Publish(testEvent{Number: i})
		}
	}()
	go func() {
		for i := 300000; i < 400000; i++ {
			eb.Publish(testEvent{Number: i})
		}
	}()
	go func() {
		for i := 400000; i < 500000; i++ {
			eb.Publish(testEvent{Number: i})
		}
	}()
	go func() {
		for i := 500000; i < 600000; i++ {
			eb.Publish(testEvent{Number: i})
		}
	}()
	go func() {
		for i := 600000; i < 700000; i++ {
			eb.Publish(testEvent{Number: i})
		}
	}()
	go func() {
		for i := 700000; i < 800000; i++ {
			eb.Publish(testEvent{Number: i})
		}
	}()
	go func() {
		for i := 800000; i < 900000; i++ {
			eb.Publish(testEvent{Number: i})
		}
	}()

	for i := 900000; i < 1000000; i++ {
		eb.Publish(testEvent{Number: i})
	}

	for counter < 1000000 {
		<-time.NewTimer(50*time.Millisecond).C
	}

	fmt.Println("done")

	cancelFunc()

}


