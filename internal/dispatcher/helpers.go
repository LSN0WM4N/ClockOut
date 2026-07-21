package dispatcher

import (
	"fmt"

	"ClockOut/internal/logger"
)

func (d *Dispatcher[T]) safeInvoke(item T) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: Notify the user that there were an error during the dispatching
			// not just a "silent" log
			logger.Print("dispatcher", "recovered from panic en handler: %v", fmt.Sprint(r))
		}
	}()
	d.handler(item)
}
