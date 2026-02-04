package notifier

import "log/slog"

func mockLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
