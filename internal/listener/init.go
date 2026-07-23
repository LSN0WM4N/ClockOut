package listener

import (
	"net/http"

	"ClockOut/internal/constants"
	"ClockOut/internal/core"
	"ClockOut/internal/logger"
)

func Init(eventsCh chan<- core.Event) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/event", newEventHandler(eventsCh))

	ln, err := resolveListener()
	if err != nil {
		return constants.ErrorInitializingServer(err)
	}

	logger.Print("listener", "Server listening at %s", ln.Addr().String())

	if err := http.Serve(ln, mux); err != nil {
		return constants.ErrorOnServer(err)
	}

	return nil
}
