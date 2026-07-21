package dispatcher

type Dispatcher[T any] struct {
	ch      <-chan T
	handler func(T)
	done    chan struct{}
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

func (d *Dispatcher[T]) Stop() {
	close(d.done)
}
