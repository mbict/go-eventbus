package eventbus

import (
	"reflect"
)

type MappedHandlers map[string][]EventHandlerFunc

// EventHandlerResolver tries to resolve all the compatible handler from a function
// or exported methods by a type instance. It will return all the found handlers in a map
// where the mapped key is the event name. When the Event interface is used as the argument type
// it will be registered as a woldcard handler "*".
//
// A function or a method in an instance should have the structure
// func( event eventbus.Event ) error
//
//	or a compatible event type that implements the evenbus.Event
//
// func( event MyEvent ) error
func EventHandlerResolver(handler any) (MappedHandlers, error) {

	v := reflect.ValueOf(handler)

	eventImpl := reflect.TypeOf((*Event)(nil)).Elem()
	errorImpl := reflect.TypeOf((*error)(nil)).Elem()

	testSuitableHandler := func(t reflect.Type) bool {
		return t.NumIn() == 1 && t.In(0).Implements(eventImpl) &&
			t.NumOut() == 1 && t.Out(0).Implements(errorImpl)
	}

	genHandler := func(v reflect.Value) func(any) error {
		eventType := v.Type().In(0)
		return func(evt any) error {
			evti := reflect.ValueOf(evt)
			if evti.CanConvert(eventType) {
				out := v.Call([]reflect.Value{evti.Convert(eventType)})
				result := out[0].Interface()
				if result != nil {
					return result.(error)
				}
				return nil
			} else {
				panic("unable to convert to eventType")
			}
		}
	}

	getEventNameFromMethod := func(v reflect.Value) string {
		if v.Type().In(0).Kind() == reflect.Interface {
			return "*"
		}
		return reflect.Zero(v.Type().In(0)).Interface().(Event).EventName()
	}

	res := make(MappedHandlers)

	//func
	if v.Kind() == reflect.Func {
		if testSuitableHandler(v.Type()) == true {
			eventName := getEventNameFromMethod(v)
			res[eventName] = append(res[eventName], genHandler(v))
		}
	}

	//methods on type
	for i := 0; i < v.NumMethod(); i++ {
		v := v.Method(i)
		if testSuitableHandler(v.Type()) {
			eventName := getEventNameFromMethod(v)
			res[eventName] = append(res[eventName], genHandler(v))
		}
	}
	return res, nil
}
