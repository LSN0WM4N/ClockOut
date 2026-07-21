package constants

import (
	"fmt"
)

// DB ERRORS

// Return a custom error based on the ID passed as argument
func EmployeeDoesNotExists(id int64) error {
	return fmt.Errorf("Employee with ID[%d] does not exists in the Database", id)
}

// SERVER ERRORS

func ErrorOnServer(err error) error {
	return fmt.Errorf("Server error: %v", err)
}

// Return the error on server Init
func ErrorInitializingServer(err error) error {
	return fmt.Errorf("Error while initializing the server: %v", err)
}

func ErrorNoAvailablePort() error {
	return fmt.Errorf("No ports for HTTP server available")
}
