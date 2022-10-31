package eventbus

import "reflect"

type EventNameResolver func(event any) string

func resolveEventName(event any) string {
	if c, ok := event.(Event); ok {
		return c.EventName()
	}

	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		return t.Elem().PkgPath() + "/*" + t.Elem().Name()
	}
	return t.PkgPath() + "/" + t.Name()
}
