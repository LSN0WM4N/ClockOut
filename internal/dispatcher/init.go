package dispatcher

import "ClockOut/internal/core"

func NewDispatcher[T core.EventInterface](ch <-chan T) *Dispatcher[T] {
	return &Dispatcher[T]{
		ch:       ch,
		handlers: make(map[string]func(T)),
	}
}
