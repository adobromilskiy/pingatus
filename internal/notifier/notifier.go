package notifier

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/adobromilskiy/pingatus/internal/config"
)

var (
	errEmptyToken  = errors.New("notifier: tgtoken is empty")
	errEmptyChatID = errors.New("notifier: tgchatid is empty")
	errUnknownType = errors.New("notifier: unknown notifier type")
)

type Notifier interface {
	Send(ctx context.Context, msg string)
}

func New(lg *slog.Logger, cfg config.NotifierConfig) (Notifier, error) {
	var notifier Notifier

	switch cfg.Type {
	case "telegram":
		if len(cfg.TgToken) == 0 {
			return nil, errEmptyToken
		}

		if len(cfg.TgChatID) == 0 {
			return nil, errEmptyChatID
		}

		notifier = newTelegram(lg, cfg.TgToken, cfg.TgChatID)

	default:
		return nil, fmt.Errorf("%w: %s", errUnknownType, cfg.Type)
	}

	return notifier, nil
}
