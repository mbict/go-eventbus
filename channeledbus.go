package eventbus

import "sync"

type CancelFunc func()

type channeledEventBus struct {
	EventBus
	c chan Event
}

func (eb *channeledEventBus) Publish(event Event) (err error) {
	eb.c <- event
	return nil
}

func NewChanneld(errorHandler PublishErrorHandlerFunc) (EventBus, CancelFunc) {
	return NewChanneldWith(NewConcurrent())
}

func NewChanneldWith(eventBus EventBus) (EventBus, CancelFunc) {
	done := make(chan bool)
	eb := &channeledEventBus{
		EventBus: eventBus,
		c:        make(chan Event, 50),
	}

	go func() {
		for true {
			select {
			case event := <-eb.c:
				_ = eb.EventBus.Publish(event)
			case <-done:
				close(eb.c)
				close(done)
				return
			}
		}
	}()

	var once sync.Once
	cancelFunc := func() {
		once.Do(func() {
			done <- true
		})
	}

	return eb, cancelFunc
}
