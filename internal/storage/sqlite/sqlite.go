package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	_ "modernc.org/sqlite" //nolint: review
)

var (
	errOpenDB  = errors.New("sqlite: failed to open database")
	errCloseDB = errors.New("sqlite: failed to close database")
)

type DB struct {
	lg  *slog.Logger
	db  *sql.DB
	dsn string
}

func New(lg *slog.Logger, dsn string) *DB {
	return &DB{
		lg:  lg.With("pkg", "sqlite"),
		dsn: strings.TrimPrefix(dsn, "sqlite://"),
	}
}

func (s *DB) Open(ctx context.Context) error {
	s.db, _ = sql.Open("sqlite", s.dsn)

	err := s.db.Ping()
	if err != nil {
		return fmt.Errorf("%w: %w", errOpenDB, err)
	}

	s.lg.Log(ctx, slog.LevelInfo, "open database", "dsn", s.dsn)

	return nil
}

func (s *DB) Close(ctx context.Context) error {
	s.lg.Log(ctx, slog.LevelInfo, "closed database", "dsn", s.dsn)

	if s.db == nil {
		return nil
	}

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%w: %w", errCloseDB, err)
	}

	return nil
}

func (s *DB) Get() *sql.DB {
	return s.db
}
