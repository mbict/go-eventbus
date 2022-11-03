package eventbus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type eventPtr struct {
}

func (e eventPtr) EventName() string {
	return "eventPtr"
}

type eventA string

func (e eventA) EventName() string {
	return "eventA"
}

type eventB string

func (e eventB) EventName() string {
	return "eventB"
}

type testHandlerInstance struct {
	called []string
}

func (t *testHandlerInstance) HandleA(eventA eventA) error {
	t.called = append(t.called, string(eventA))
	return nil
}

func (t *testHandlerInstance) HandleB(eventB eventB) error {
	t.called = append(t.called, string(eventB))
	return nil
}

func (t *testHandlerInstance) HandleCatchAll(eventB Event) error {
	t.called = append(t.called, eventB.EventName())
	return nil
}

func (t *testHandlerInstance) HandleIgnore(eventC int) error {
	t.called = append(t.called, "error")
	return nil
}

func Test_RegisterHandlerWithInstanceMethods(t *testing.T) {

	h := &testHandlerInstance{}

	res, err := EventHandlerResolver(h)
	assert.NoError(t, err)
	assert.Len(t, res, 3)

	assert.Contains(t, res, "eventA")
	assert.Len(t, res["eventA"], 1)
	assert.Contains(t, res, "eventB")
	assert.Len(t, res["eventB"], 1)
	assert.Contains(t, res, "*")
	assert.Len(t, res["*"], 1)

	//test calling the handlers
	err = res["eventA"][0](eventA("1"))
	assert.NoError(t, err)

	err = res["eventB"][0](eventA("2"))
	assert.NoError(t, err)

	err = res["*"][0](eventA("this would be the event name"))
	assert.NoError(t, err)

	assert.Equal(t, []string{"1", "2", "eventA"}, h.called)
}

func Test_RegisterHandlerWithFunc(t *testing.T) {

	var called []string

	h := func(event eventA) error { called = append(called, string(event)); return nil }

	res, err := EventHandlerResolver(h)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Contains(t, res, "eventA")
	assert.Len(t, res["eventA"], 1)

	//calling the returned method
	err = res["eventA"][0](eventA("test"))
	assert.NoError(t, err)
	assert.Len(t, called, 1)
	assert.Equal(t, []string{"test"}, called)
}

func Test_RegisterHandlerWithFuncAsWildcard(t *testing.T) {
	called := []string{}

	// As the handler only specifies the evenbus.Event interface, we cannot determine the event name, and we assume the handler is intended for wildcard usage
	h := func(event Event) error { called = append(called, string(event.(eventA))); return nil }

	res, err := EventHandlerResolver(h)

	assert.NoError(t, err)
	assert.Len(t, res, 1)

	assert.Contains(t, res, "*")
	assert.Len(t, res["*"], 1)

	//calling the returned method
	err = res["*"][0](eventA("test"))
	assert.NoError(t, err)
	assert.Len(t, called, 1)
	assert.Equal(t, []string{"test"}, called)
}

func Test_RegisterHandlerWithFunc_non_compatible(t *testing.T) {

	handlerCalled := 0

	//the handler function first argument does not implement the eventbus.Event, thus we ignore the handler
	h := func(event int) error { handlerCalled += 1; return nil }

	res, err := EventHandlerResolver(h)

	assert.NoError(t, err)
	assert.Len(t, res, 0)
}

func Test_RegisterHandlerWithPointerEventAndNonPointerReceiverConversion(t *testing.T) {
	called := []string{}

	// As the handler only specifies the evenbus.Event interface, we cannot determine the event name, and we assume the handler is intended for wildcard usage
	h := func(event eventPtr) error { called = append(called, event.EventName()); return nil }

	res, err := EventHandlerResolver(h)

	assert.NoError(t, err)
	assert.Len(t, res, 1)

	assert.Contains(t, res, "eventPtr")
	assert.Len(t, res["eventPtr"], 1)

	//calling the returned method
	err = res["eventPtr"][0](&eventPtr{})
	assert.NoError(t, err)
	assert.Len(t, called, 1)
	assert.Equal(t, []string{"eventPtr"}, called)
}
