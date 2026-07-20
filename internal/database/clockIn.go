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
