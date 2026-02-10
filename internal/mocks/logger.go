package mocks

import "log/slog"

func MockLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
