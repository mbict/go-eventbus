package _bench

import (
	"github.com/mbict/go-eventbus/v2"
	"sync"
	"testing"
)

type testEvent struct {
	Number int
}

func (testEvent) EventName() eventbus.EventName {
	return "test.event"
}

func eventHandler(event any) error {
	e := event.(testEvent)
	for i := 0; i <= 1000; i++ {
		e.Number += i
	}
	return nil
}

func BenchmarkSimpleBus(b *testing.B) {
	eb := eventbus.New()
	eb.Subscribe(eventbus.EventHandlerFunc(eventHandler))
	e := testEvent{}

	publishEventsLoop(b, eb, e)
}

func BenchmarkConcurrentBus(b *testing.B) {
	eb := eventbus.NewConcurrent()
	eb.Subscribe(eventbus.EventHandlerFunc(eventHandler))
	e := testEvent{}

	publishEventsLoop(b, eb, e)
}

func BenchmarkAsyncBus(b *testing.B) {
	eb := eventbus.NewAsync()
	eb.Subscribe(eventbus.EventHandlerFunc(eventHandler))
	e := testEvent{}

	publishEventsLoop(b, eb, e)
}

func BenchmarkChanneldWithSimpleBus(b *testing.B) {

	eb, cancelFunc := eventbus.NewChanneldWith(eventbus.New())
	eb.Subscribe(eventbus.EventHandlerFunc(eventHandler))
	e := testEvent{}

	publishEventsLoop(b, eb, e)

	cancelFunc()
}

func BenchmarkChanneldWithCuncurrentBus(b *testing.B) {
	eb, cancelFunc := eventbus.NewChanneldWith(eventbus.NewConcurrent())
	eb.Subscribe(eventbus.EventHandlerFunc(eventHandler))
	e := testEvent{}

	publishEventsLoop(b, eb, e)

	cancelFunc()
}

func BenchmarkChanneldWithAsyncBus(b *testing.B) {
	eb, cancelFunc := eventbus.NewChanneldWith(eventbus.NewAsync())
	eb.Subscribe(eventbus.EventHandlerFunc(eventHandler))
	e := testEvent{}

	publishEventsLoop(b, eb, e)

	cancelFunc()
}

func BenchmarkPlainChannel(b *testing.B) {
	c := make(chan eventbus.Event)
	done := make(chan bool)
	e := testEvent{}

	go func() {
		for true {
			select {
			case evt := <-c:
				eventHandler(evt)
			case <-done:
				return
			}
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < b.N; n++ {
		for pp := 0; pp < 10; pp++ {
			wg.Add(1)
			go func() {
				for i := 0; i < 100; i++ {
					c <- e
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}

	done <- true
}

func publishEventsLoop(b *testing.B, eb eventbus.EventBus, e testEvent) {
	var wg sync.WaitGroup
	for n := 0; n < b.N; n++ {
		for pp := 0; pp < 10; pp++ {
			wg.Add(1)
			go func() {
				for i := 0; i < 100; i++ {
					eb.Publish(e)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
