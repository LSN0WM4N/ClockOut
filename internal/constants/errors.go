package constants

import (
	"fmt"
)

// DB ERRORS

// Return a custom error based on the ID passed as argument
func EmployeeDoesNotExists(id int64) error {
	return fmt.Errorf("Employee with ID[%d] does not exists in the Database", id)
}
