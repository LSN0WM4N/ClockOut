package database

import (
	"database/sql"
	"time"

	"ClockOut/internal/constants"
	"ClockOut/internal/model"
)

// Insert a new ClockIn entry in the database for the employee `employeeId`
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
		EmployeeId: employeeID,
		Timestamp:  now,
		Type:       t,
	}, nil
}

// Return the `limit` newer entries in the clock_in database,
// skip the first`offset` entries
// Return an array of *model.ClockIn or an error
func GetAllClockIn(db *sql.DB, offset, limit int) ([]*model.ClockIn, error) {
	// No limit over the query
	if limit <= 0 {
		limit = -1
	}

	rows, err := db.Query(`
		SELECT
			id,
			employee_id,
			timestamp,
			type
		FROM clock_in
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?
	`, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	clockIns, err := scanClockIn(rows)

	if err != nil {
		return nil, err
	}

	return clockIns, nil
}

// Get the last `limit` clock_in entries for a given employee (if
// limit = -1 return all of them) skipping the last 'offset  ones
func GetClockInsByEmployee(db *sql.DB, employeeId int64, offset int, limit int) ([]*model.ClockIn, error) {
	if limit <= 0 {
		limit = -1
	}

	rows, err := db.Query(`
		SELECT
			id,
			employee_id,
			timestamp,
			type
		FROM clock_in
		WHERE employee_id = ?
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?
	`, employeeId, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	clockIns, err := scanClockIn(rows)

	if err != nil {
		return nil, err
	}

	return clockIns, nil
}
