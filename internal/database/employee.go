package database

import (
	"ClockOut/internal/model"
	"database/sql"
	"strconv"
	"time"
)

func InsertEmployee(db *sql.DB, name, role string) (*model.Employee, error) {
	now := time.Now()

	result, err := db.Exec(`
		INSERT INTO employee (
			name,
			role,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?)`,
		name,
		role,
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

	return &model.Employee{
		ID:   strconv.FormatInt(id, 10),
		Name: name,
		Role: role,
	}, nil
}

// Return if the employee is working (last operation was a clock in [start])
// or if is resting (last operation was a clock out [finish])
func CheckEmployeeStatus(db *sql.DB, id int64) (*model.Type, error) {
	var status model.Type

	err := db.QueryRow(`
		SELECT type
		FROM clock_in
		WHERE employee_id = ?
		ORDER BY timestamp DESC
		LIMIT 1
	`, id).Scan(&status)

	if err != nil {
		// Due to practical effects, if there is no Clock In, i will
		// take the next assignation as Clock Out
		if err == sql.ErrNoRows {
			status = model.Finish
			return &status, nil
		}
		return nil, err
	}

	return &status, nil
}

// Return True or False if the queried employee exists or not
// In case of an error, the result will be false
func CheckEmployeeExists(db *sql.DB, id int64) (bool, error) {
	var exists bool

	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM employee
			WHERE id = ?
		)
	`, id).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
