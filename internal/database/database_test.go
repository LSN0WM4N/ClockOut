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

func TestInsertNoNameEmployee(t *testing.T) {
	db := setupDB(t)

	_, err := database.InsertEmployee(db, "", "Developer")
	if err == nil {
		t.Fatal("expected error when inserting employee with empty name, got nil")
	}
}

func TestInsertNoRoleEmployee(t *testing.T) {
	db := setupDB(t)

	_, err := database.InsertEmployee(db, "John", "")
	if err == nil {
		t.Fatal("expected error when inserting employee with empty role, got nil")
	}
}

func TestInsertInvalidEmployeeClockIn(t *testing.T) {
	db := setupDB(t)

	// Insert employee first
	employee, err := database.InsertEmployee(db, "John", "Developer")
	if err != nil {
		t.Fatalf("insert employee: %v", err)
	}

	// Try to insert clock-in with invalid employee ID
	invalidID := employee.ID + 999
	_, err = database.InsertClockIn(db, invalidID, model.Start)
	if err == nil {
		t.Fatal("expected error when inserting clock-in with invalid employee ID, got nil")
	}
}

func TestInsertInvalidTypeClockIn(t *testing.T) {
	db := setupDB(t)

	// Insert employee first
	employee, err := database.InsertEmployee(db, "John", "Developer")
	if err != nil {
		t.Fatalf("insert employee: %v", err)
	}

	// Try to insert clock-in with invalid type
	_, err = database.InsertClockIn(db, employee.ID, "invalid_type")
	if err == nil {
		t.Fatal("expected error when inserting clock-in with invalid type, got nil")
	}
}

func TestInsertClockIn(t *testing.T) {
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

// Additional tests for better coverage

func TestGetEmployeeByID(t *testing.T) {
	db := setupDB(t)

	inserted, err := database.InsertEmployee(db, "John", "Developer")
	if err != nil {
		t.Fatalf("insert employee: %v", err)
	}

	employee, err := database.GetEmployeeByID(db, inserted.ID)
	if err != nil {
		t.Fatalf("get employee by ID: %v", err)
	}

	if employee.ID != inserted.ID {
		t.Fatalf("expected ID %d, got %d", inserted.ID, employee.ID)
	}
	if employee.Name != "John" {
		t.Fatalf("expected name John, got %s", employee.Name)
	}
	if employee.Role != "Developer" {
		t.Fatalf("expected role Developer, got %s", employee.Role)
	}
}

func TestGetEmployeeByIDNotFound(t *testing.T) {
	db := setupDB(t)

	_, err := database.GetEmployeeByID(db, 999)
	if err == nil {
		t.Fatal("expected error when getting non-existent employee, got nil")
	}
}

func TestGetAllEmployees(t *testing.T) {
	db := setupDB(t)

	// Insert multiple employees
	employees := []struct {
		name string
		role string
	}{
		{"Alice", "Manager"},
		{"Bob", "Developer"},
		{"Charlie", "Designer"},
	}

	for _, e := range employees {
		_, err := database.InsertEmployee(db, e.name, e.role)
		if err != nil {
			t.Fatalf("insert employee: %v", err)
		}
	}

	allEmployees, err := database.GetAllEmployees(db, 0, 0)
	if err != nil {
		t.Fatalf("get all employees: %v", err)
	}

	if len(allEmployees) != len(employees) {
		t.Fatalf("expected %d employees, got %d", len(employees), len(allEmployees))
	}
}

func TestGetClockInsByEmployee(t *testing.T) {
	db := setupDB(t)

	employee, err := database.InsertEmployee(db, "John", "Developer")
	if err != nil {
		t.Fatalf("insert employee: %v", err)
	}

	// Insert multiple clock-ins
	for i := 0; i < 3; i++ {
		_, err := database.InsertClockIn(db, employee.ID, model.Start)
		if err != nil {
			t.Fatalf("insert clock-in: %v", err)
		}
	}

	clockIns, err := database.GetClockInsByEmployee(db, employee.ID, 0, 0)
	if err != nil {
		t.Fatalf("get clock-ins by employee: %v", err)
	}

	if len(clockIns) != 3 {
		t.Fatalf("expected 3 clock-ins, got %d", len(clockIns))
	}
}

func TestDeleteEmployee(t *testing.T) {
	db := setupDB(t)

	employee, err := database.InsertEmployee(db, "John", "Developer")
	if err != nil {
		t.Fatalf("insert employee: %v", err)
	}

	err = database.DeleteEmployee(db, employee.ID)
	if err != nil {
		t.Fatalf("delete employee: %v", err)
	}

	// Verify employee is deleted
	_, err = database.GetEmployeeByID(db, employee.ID)
	if err == nil {
		t.Fatal("expected error when getting deleted employee, got nil")
	}
}

func TestUpdateEmployee(t *testing.T) {
	db := setupDB(t)

	employee, err := database.InsertEmployee(db, "John", "Developer")
	if err != nil {
		t.Fatalf("insert employee: %v", err)
	}

	newEmployee := model.Employee{
		ID:   employee.ID,
		Name: "John Updated",
		Role: "Senior Developer",
	}

	err = database.UpdateEmployee(db, &newEmployee)
	if err != nil {
		t.Fatalf("update employee: %v", err)
	}

	updated, err := database.GetEmployeeByID(db, employee.ID)
	if err != nil {
		t.Fatalf("get updated employee: %v", err)
	}

	if updated.Name != "John Updated" {
		t.Fatalf("expected name John Updated, got %s", updated.Name)
	}
	if updated.Role != "Senior Developer" {
		t.Fatalf("expected role Senior Developer, got %s", updated.Role)
	}
}
