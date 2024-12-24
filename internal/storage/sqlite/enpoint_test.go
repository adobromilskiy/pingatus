package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/adobromilskiy/pingatus/core"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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

	err = endpoint.Save(context.TODO(), data)
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

func TestEndpoint_Save_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	mock.ExpectExec("INSERT INTO endpoints").
		WithArgs("test_name", "test_address", false, 1672531200).
		WillReturnError(errors.New("mocked database error"))

	e := &Endpoint{db: db}

	data := core.Endpoint{
		Name:    "test_name",
		Address: "test_address",
		Status:  false,
		Date:    1672531200,
	}

	err = e.Save(context.TODO(), data)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetEndpoints(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	e := &Endpoint{db: db}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("endpoint1").
			AddRow("endpoint2")
		mock.ExpectQuery(`SELECT DISTINCT name FROM endpoints;`).WillReturnRows(rows)

		result, err := e.GetEndpoints(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, []string{"endpoint1", "endpoint2"}, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT DISTINCT name FROM endpoints;`).WillReturnError(sql.ErrConnDone)

		result, err := e.GetEndpoints(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("endpoint1").
			AddRow(nil)
		mock.ExpectQuery(`SELECT DISTINCT name FROM endpoints;`).WillReturnRows(rows)

		result, err := e.GetEndpoints(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("endpoint1").
			RowError(0, sql.ErrNoRows)
		mock.ExpectQuery(`SELECT DISTINCT name FROM endpoints;`).WillReturnRows(rows)

		result, err := e.GetEndpoints(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetEndpointStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	e := &Endpoint{db: db}

	const query = `SELECT name, address, status, date FROM endpoints WHERE name=\? and date > \? ORDER BY date ASC;`

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name", "address", "status", "date"}).
			AddRow("endpoint1", "192.168.1.1", true, 1660000000).
			AddRow("endpoint1", "192.168.1.1", false, 1660000100)

		mock.ExpectQuery(query).WithArgs("endpoint1", int64(1659999999)).WillReturnRows(rows)

		result, err := e.GetEndpointStats(context.Background(), "endpoint1", 1659999999)

		assert.NoError(t, err)
		expected := []core.Endpoint{
			{Name: "endpoint1", Address: "192.168.1.1", Status: true, Date: 1660000000},
			{Name: "endpoint1", Address: "192.168.1.1", Status: false, Date: 1660000100},
		}
		assert.Equal(t, expected, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(query).WithArgs("endpoint1", int64(1659999999)).WillReturnError(fmt.Errorf("query failed"))

		result, err := e.GetEndpointStats(context.Background(), "endpoint1", 1659999999)

		assert.Error(t, err)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name", "address", "status", "date"}).
			AddRow("endpoint1", nil, true, 1660000000)

		mock.ExpectQuery(query).WithArgs("endpoint1", int64(1659999999)).WillReturnRows(rows)

		result, err := e.GetEndpointStats(context.Background(), "endpoint1", 1659999999)

		assert.Error(t, err)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name", "address", "status", "date"}).
			AddRow("endpoint1", "192.168.1.1", true, 1660000000).
			RowError(0, fmt.Errorf("row error"))

		mock.ExpectQuery(query).WithArgs("endpoint1", int64(1659999999)).WillReturnRows(rows)

		result, err := e.GetEndpointStats(context.Background(), "endpoint1", 1659999999)

		assert.Error(t, err)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
