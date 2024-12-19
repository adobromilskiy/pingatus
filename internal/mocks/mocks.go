package mocks

import (
	"io"
	"log/slog"
)

func MockLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
