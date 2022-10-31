package eventbus

import "sync"

type CancelFunc func()

type channeledEventBus struct {
	EventBus
	c chan any
}

func (eb *channeledEventBus) Publish(event any) (err error) {
	eb.c <- event
	return nil
}

func NewChanneldWith(eventBus EventBus) (EventBus, CancelFunc) {
	done := make(chan bool)
	eb := &channeledEventBus{
		EventBus: eventBus,
		c:        make(chan any, 100),
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
