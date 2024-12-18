package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/adobromilskiy/pingatus/core"
)

type Endpoint struct {
	db *sql.DB
}

func NewEndpoint(db *sql.DB) (*Endpoint, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS endpoints (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			status BOOLEAN NOT NULL,
			date INTEGER NOT NULL
		);`)
	if err != nil {
		return nil, fmt.Errorf("endpoint: can not create table: %w", err)
	}

	return &Endpoint{
		db: db,
	}, nil
}

func (e *Endpoint) Save(data core.Endpoint) error {
	const query = `
		INSERT INTO endpoints
		(name, address, status, date)
		VALUES (?, ?, ?, ?);
	`

	_, err := e.db.Exec(query, data.Name, data.Address, data.Status, data.Date)
	if err != nil {
		return fmt.Errorf("endpoint: failed to save: %w", err)
	}

	return nil
}
