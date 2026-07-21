package listener

import (
	"net/http"

	"ClockOut/internal/constants"
	"ClockOut/internal/logger"
)

func Init(eventsCh chan<- string) error {
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
