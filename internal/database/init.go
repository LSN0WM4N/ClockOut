package database

import "database/sql"

func Init(db *sql.DB) error {

	query := `
		CREATE TABLE IF NOT EXISTS employee (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			role TEXT,

			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS clock_in (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			employee_id INTEGER NOT NULL,
			timestamp DATETIME NOT NULL,
			type TEXT NOT NULL CHECK(type IN ('entry', 'exit')),

			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,

			FOREIGN KEY (employee_id) REFERENCES employee(id)
		);
	`

	_, err := db.Exec(query)
	return err
}
