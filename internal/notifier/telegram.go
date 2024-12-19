package notifier

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

type telegram struct {
	lg     *slog.Logger
	token  string
	chatID string
	apiURL string
}

func newTelegram(lg *slog.Logger, token, chatID string) *telegram {
	lg = lg.With("pkg", "notifier")

	lg.Log(context.Background(), slog.LevelInfo, "telegram notifier initialized", "token", token, "chat_id", chatID)

	return &telegram{lg, token, chatID, "https://api.telegram.org/bot" + token}
}

func (t *telegram) Send(ctx context.Context, msg string) {
	data := url.Values{
		"chat_id": {t.chatID},
		"text":    {msg},
	}

	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/sendMessage?%s", t.apiURL, data.Encode()), nil)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.lg.Log(ctx, slog.LevelError, "failed to send notificiation via telegram", "err", err)

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.lg.Log(ctx, slog.LevelError, "failed to send notification via telegram", "status", resp.Status)
	}
}
