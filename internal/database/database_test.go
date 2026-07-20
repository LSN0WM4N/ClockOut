package database_test

import (
	"database/sql"
	"os"
	"testing"

	"ClockOut/internal/database"
	"ClockOut/internal/model"

	_ "modernc.org/sqlite"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()

	const dbPath = "test.db"

	_ = os.Remove(dbPath)

	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := database.Init(db); err != nil {
		t.Fatalf("init db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(dbPath)
	})

	return db
}

func TestInsertEmployee(t *testing.T) {
	db := setupDB(t)

	_, err := database.InsertEmployee(db, "John", "Developer")
	if err != nil {
		t.Fatalf("insert employee: %v", err)
	}
}

func testInsertNoNameEmployee()
func testInsertNoRoleEmployee()
func testInsertInvalidEmployeeClockIn()
func testInsertInvalidTypeClockIn()

func testInsertClockIn(t *testing.T) {
	db := setupDB(t)

	employee, err := database.InsertEmployee(db, "John", "Developer")
	if err != nil {
		t.Fatalf("insert employee: %v", err)
	}

	_, err = database.InsertClockIn(db, employee.ID, model.Start)
	if err != nil {
		t.Fatalf("inserting clock_in: %v", err)
	}
}

func TestClockInFlow(t *testing.T) {
	db := setupDB(t)

	employee, err := database.InsertEmployee(db, "Test", "Employee")
	if err != nil {
		t.Fatal(err)
	}

	status, exists, err := database.CheckEmployeeStatus(db, employee.ID)
	if err != nil {
		t.Fatal(err)
	}

	if exists {
		t.Fatal("employee should have no status at this point")
	}

	_, err = database.InsertClockIn(db, employee.ID, model.Start)
	if err != nil {
		t.Fatal(err)
	}

	status, exists, err = database.CheckEmployeeStatus(db, employee.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("employee should have a status")
	}

	if status != model.Start {
		t.Fatalf("expected %s, got %s", model.Start, status)
	}

	_, err = database.InsertClockIn(db, employee.ID, model.Finish)
	if err != nil {
		t.Fatal(err)
	}

	status, exists, err = database.CheckEmployeeStatus(db, employee.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("employee should have a status")
	}

	if status != model.Finish {
		t.Fatalf("expected %s, got %s", model.Finish, status)
	}

	entries, err := database.GetAllClockIn(db, 0, -1)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}
