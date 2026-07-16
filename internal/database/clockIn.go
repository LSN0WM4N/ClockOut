package database

import (
	"database/sql"
	"strconv"
	"time"

	"ClockOut/internal/constants"
	"ClockOut/internal/model"
)

func InsertClockIn(db *sql.DB, employeeID int64, t model.Type) (*model.ClockIn, error) {
	now := time.Now()

	if employeeExists, err := CheckEmployeeExists(db, employeeID); err != nil || !employeeExists {
		return nil, constants.EmployeeDoesNotExists(employeeID)
	}

	result, err := db.Exec(`
		INSERT INTO clock_in (
			employee_id,
			timestamp,
			type,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?)`,
		employeeID,
		now,
		t,
		now,
		now,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.ClockIn{
		ID:         id,
		EmployeeId: strconv.FormatInt(employeeID, 10),
		Timestamp:  now,
		Type:       t,
	}, nil
}
