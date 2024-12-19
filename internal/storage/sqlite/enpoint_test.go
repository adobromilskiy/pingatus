package sqlite

import (
	"database/sql"
	"testing"

	"github.com/adobromilskiy/pingatus/core"
)

func TestNewEndpoint_Success(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = NewEndpoint(db)
	if err != nil {
		t.Fatalf("NewEndpoint() failed: %v", err)
	}

	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='endpoints';").Scan(&tableName)
	if err != nil {
		t.Fatalf("Table 'endpoints' not found: %v", err)
	}
	if tableName != "endpoints" {
		t.Fatalf("Expected table name 'endpoints', got '%s'", tableName)
	}
}

func TestEndpoint_Save_Success(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	endpoint, err := NewEndpoint(db)
	if err != nil {
		t.Fatalf("NewEndpoint() failed: %v", err)
	}

	data := core.Endpoint{
		Name:    "test-endpoint",
		Address: "http://localhost",
		Status:  true,
		Date:    1672531200,
	}

	err = endpoint.Save(data)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	var (
		name    string
		address string
		status  bool
		date    int64
	)

	err = db.QueryRow("SELECT name, address, status, date FROM endpoints WHERE name=?;", data.Name).Scan(&name, &address, &status, &date)
	if err != nil {
		t.Fatalf("QueryRow() failed: %v", err)
	}

	if name != data.Name {
		t.Errorf("Expected name %s, got %s", data.Name, name)
	}

	if address != data.Address {
		t.Errorf("Expected address %s, got %s", data.Address, address)
	}

	if status != data.Status {
		t.Errorf("Expected status %t, got %t", data.Status, status)
	}

	if date != data.Date {
		t.Errorf("Expected date %d, got %d", data.Date, date)
	}
}
