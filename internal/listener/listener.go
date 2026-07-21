package listener

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"ClockOut/internal/constants"
)

func resolveListener() (net.Listener, error) {
	// Try Env
	if envPort := os.Getenv(EnvPortName); envPort != "" {
		if port, err := strconv.Atoi(envPort); err == nil && port > 0 && port <= 65535 {
			addr := fmt.Sprintf(":%d", port)
			if ln, err := net.Listen("tcp", addr); err == nil {
				return ln, nil
			}
		}
	}

	addr := fmt.Sprintf(":%d", DefaultPort)
	if ln, err := net.Listen("tcp", addr); err == nil {
		return ln, nil
	}

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, constants.ErrorNoAvailablePort()
	}

	return ln, nil
}
