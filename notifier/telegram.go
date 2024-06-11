package notifier

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Telegram struct {
	Token  string
	ChatID string
}

func NewTelegram(token, chatID string) *Telegram {
	return &Telegram{token, chatID}
}

func (t *Telegram) Send(msg string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)
	data := url.Values{
		"chat_id": {t.ChatID},
		"text":    {msg},
	}

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		log.Printf("[ERROR] failed to send notificiation via telegram: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] failed to send notification via telegram: %v", resp.Status)
	}
}
