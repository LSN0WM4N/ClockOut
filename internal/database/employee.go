package database

import (
	"ClockOut/internal/model"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func InsertEmployee(db *sql.DB, name, role string) (*model.Employee, error) {
	now := time.Now()

	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("Invalid name")
	}

	if strings.TrimSpace(role) == "" {
		return nil, fmt.Errorf("Invalid role")
	}

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
		ID:   id,
		Name: name,
		Role: role,
	}, nil
}

// Return if the employee is working (last operation was a clock in [start])
// or if is resting (last operation was a clock out [finish]) and a boolean
// that indicates that there are no entries for that employee
func CheckEmployeeStatus(db *sql.DB, id int64) (model.Type, bool, error) {
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
			return model.Finish, false, nil
		}
		return model.Finish, false, err
	}

	return status, true, nil
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

// Return the employee with the given ID
func GetEmployeeByID(db *sql.DB, id int64) (*model.Employee, error) {
	var employee model.Employee

	err := db.QueryRow(` 
		SELECT 
			id, 
			name,
			role
		FROM employee
		WHERE id = ? 
	`, id).Scan(
		&employee.ID,
		&employee.Name,
		&employee.Role,
	)

	if err != nil {
		return nil, err
	}
	return &employee, err
}

// Return the first `limit` employees (if limit = -1 then return all of them)
// skipping the first `offset` ones
func GetAllEmployees(db *sql.DB, offset, limit int) ([]*model.Employee, error) {
	if limit <= 0 {
		limit = -1
	}

	rows, err := db.Query(`
		SELECT
			id,
			name,
			role
		FROM employee
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*model.Employee

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

		employees = append(employees, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

// Delete the given employee from the Database
func DeleteEmployee(db *sql.DB, id int64) error {
	rows, err := db.Exec(`
		DELETE
			FROM employee
			WHERE id = ? 
	`, id)

	if err != nil {
		return err
	}

	if affected, err := rows.RowsAffected(); affected == 0 || err != nil {
		return err
	}

	return nil
}

func UpdateEmployee(db *sql.DB, employee *model.Employee) error {
	result, err := db.Exec(`
		UPDATE employee
		SET
			name = ?,
			role = ?,
			updated_at = ?
		WHERE id = ?
	`,
		employee.Name,
		employee.Role,
		time.Now(),
		employee.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
