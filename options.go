package eventbus

type eventNameResolverSetter interface {
	setEventResolver(EventNameResolver)
}

type errorHandlerSetter interface {
	setErrorHandler(errorHandler PublishErrorHandlerFunc)
}

type Option func(EventBus)

func WithEventNameResolver(resolver EventNameResolver) Option {
	return func(bus EventBus) {
		bus.(eventNameResolverSetter).setEventResolver(resolver)
	}
}

func WithErrorHandler(errorHandler PublishErrorHandlerFunc) Option {
	return func(bus EventBus) {
		bus.(errorHandlerSetter).setErrorHandler(errorHandler)
	}
}
