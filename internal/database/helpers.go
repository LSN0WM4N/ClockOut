package database

import (
	"ClockOut/internal/model"
	"database/sql"
	"fmt"
	"strings"
)

// Makes the scan over an Employees sql query
func scanEmployee(rows *sql.Rows) ([]*model.Employee, error) {
	var Employees []*model.Employee

	for rows.Next() {
		var e model.Employee

		err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Role,
		)
		if err != nil {
			return nil, err
		}

		Employees = append(Employees, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return Employees, nil
}

// Makes the scan over an ClockIn sql query
func scanClockIn(rows *sql.Rows) ([]*model.ClockIn, error) {
	var clockIns []*model.ClockIn

	for rows.Next() {
		var c model.ClockIn

		err := rows.Scan(
			&c.ID,
			&c.EmployeeId,
			&c.Timestamp,
			&c.Type,
		)
		if err != nil {
			return nil, err
		}

		clockIns = append(clockIns, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clockIns, nil
}

// Return an error if the given name is not valid
func validateEmployeeName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("Empty string is not a valid name")
	}
	return nil
}

// Return an error if the given role is not valid
//
// This is just a very basic role checkin, in the future maybe I will
// add a valid role list or something like that
func validateEmployeeRole(role string) error {
	if strings.TrimSpace(role) == "" {
		return fmt.Errorf("Empty string is not a valid role")
	}
	return nil
}
