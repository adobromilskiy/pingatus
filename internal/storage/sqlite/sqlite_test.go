package sqlite

import (
	"errors"
	"os"
	"testing"

	"github.com/adobromilskiy/pingatus/internal/mocks"
)

func TestNew(t *testing.T) {
	logger := mocks.MockLogger()
	dsn := "sqlite://./test.db"

	sqliteDB := New(logger, dsn)

	if sqliteDB == nil {
		t.Fatal("New() returned nil")
	}

	expectedDSN := "./test.db"
	if sqliteDB.dsn != expectedDSN {
		t.Errorf("Expected DSN %s, got %s", expectedDSN, sqliteDB.dsn)
	}
}

func TestOpenAndClose_Success(t *testing.T) {
	ctx := t.Context()
	logger := mocks.MockLogger()

	tempDB := "./test_open_close.db"
	defer os.Remove(tempDB)

	sqliteDB := New(logger, tempDB)

	err := sqliteDB.Open(ctx)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}

	err = sqliteDB.Close(ctx)
	if err != nil {
		t.Fatalf("Close() failed: %v", err)
	}
}

func TestOpen_Fail(t *testing.T) {
	ctx := t.Context()
	logger := mocks.MockLogger()

	sqliteDB := New(logger, "/invalid/path/to/db.sqlite")

	err := sqliteDB.Open(ctx)
	if err == nil {
		t.Fatal("Expected error, got nil for invalid DB path")
	}

	if !errors.Is(err, errOpenDB) {
		t.Errorf("Expected error %v, got %v", errOpenDB, err)
	}
}

func TestClose_Fail(t *testing.T) {
	ctx := t.Context()
	logger := mocks.MockLogger()

	sqliteDB := &DB{
		lg: logger,
		db: nil,
	}

	err := sqliteDB.Close(ctx)
	if err != nil {
		t.Fatal("Expected no error when closing nil db")
	}
}
