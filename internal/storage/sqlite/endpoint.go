package sqlite

import (
	"context"
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

func (e *Endpoint) Save(ctx context.Context, data core.Endpoint) error {
	const query = `
		INSERT INTO endpoints
		(name, address, status, date)
		VALUES (?, ?, ?, ?);
	`

	_, err := e.db.ExecContext(ctx, query, data.Name, data.Address, data.Status, data.Date)
	if err != nil {
		return fmt.Errorf("endpoint: failed to save: %w", err)
	}

	return nil
}

func (e *Endpoint) GetEndpoints(ctx context.Context) ([]string, error) {
	const query = `SELECT DISTINCT name FROM endpoints;`

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("endpoint: failed to get endpoints: %w", err)
	}

	defer rows.Close()

	var names []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("endpoint: failed to scan: %w", err)
		}

		names = append(names, name)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("endpoint: %w", err)
	}

	return names, nil
}

func (e *Endpoint) GetEndpointStats(ctx context.Context, name string, date int64) ([]core.Endpoint, error) {
	const query = `SELECT name, address, status, date FROM endpoints WHERE name=? and date > ? ORDER BY date ASC;`

	rows, err := e.db.QueryContext(ctx, query, name, date)
	if err != nil {
		return nil, fmt.Errorf("endpoint: failed to get endpoints stats: %w", err)
	}

	defer rows.Close()

	var data []core.Endpoint

	for rows.Next() {
		var point core.Endpoint
		if err := rows.Scan(&point.Name, &point.Address, &point.Status, &point.Date); err != nil {
			return nil, fmt.Errorf("endpoint: failed to scan: %w", err)
		}

		data = append(data, point)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("endpoint: %w", err)
	}

	return data, nil
}

func (e *Endpoint) GetLastSuccess(ctx context.Context, name string) (*core.Endpoint, error) {
	const query = `SELECT name, address, status, date FROM endpoints WHERE name=? AND status=1 ORDER BY date DESC LIMIT 1;`

	row := e.db.QueryRowContext(ctx, query, name)

	var data core.Endpoint
	if err := row.Scan(&data.Name, &data.Address, &data.Status, &data.Date); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("endpoint: failed to get last success: %w", err)
	}

	return &data, nil
}

func (e *Endpoint) GetLastFailure(ctx context.Context, name string) (*core.Endpoint, error) {
	const query = `SELECT name, address, status, date FROM endpoints WHERE name=? AND status=0 ORDER BY date DESC LIMIT 1;`

	row := e.db.QueryRowContext(ctx, query, name)

	var data core.Endpoint
	if err := row.Scan(&data.Name, &data.Address, &data.Status, &data.Date); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("endpoint: failed to get last failure: %w", err)
	}

	return &data, nil
}
