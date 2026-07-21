package dispatcher

func NewDispatcher[T any](ch <-chan T, handler func(T)) *Dispatcher[T] {
	return &Dispatcher[T]{
		ch:      ch,
		handler: handler,
	}
}
