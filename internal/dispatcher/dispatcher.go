package dispatcher

import (
	"ClockOut/internal/core"
	"ClockOut/internal/logger"
	"fmt"
)

// TODO: Make the dispatcher strong against race conditions

type Dispatcher[T core.EventInterface] struct {
	ch       <-chan T
	handlers map[string]func(T)
	done     chan struct{}
}

func (d *Dispatcher[T]) Start() {
	go func() {
		for {
			select {
			case item, ok := <-d.ch:
				if !ok {
					return
				}
				go d.safeInvoke(item)
			case <-d.done:
				return
			}
		}
	}()
}

func (d *Dispatcher[T]) RegisterHandler(s string, handler func(T)) {
	if d.handlers[s] != nil {
		logger.Debug("Dispatcher", "Overwriting handler for event %s", s)
	}

	d.handlers[s] = handler
	logger.Print("dispatcher", "Register handler for event %s", s)
}

func (d *Dispatcher[T]) Stop() {
	close(d.done)
}

// Private

func (d *Dispatcher[T]) safeInvoke(item T) {
	event := T.GetType(item)

	defer func() {
		if r := recover(); r != nil {
			// TODO: Notify the user that there were an error during the dispatching
			// not just a "silent" log
			logger.Print("dispatcher", "recovered from panic en handler: ", fmt.Sprint(r))
		}
	}()
	handler, ok := d.handlers[event]
	if !ok || handler == nil {
		logger.Print("dispatcher", "no handler registered for event ", event)
		return
	}

	handler(item)
}
