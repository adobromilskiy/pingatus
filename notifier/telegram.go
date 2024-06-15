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
	APIURL string
}

func NewTelegram(token, chatID string) *Telegram {
	return &Telegram{token, chatID, fmt.Sprintf("https://api.telegram.org/bot%s", token)}
}

func (t *Telegram) Send(msg string) {
	data := url.Values{
		"chat_id": {t.ChatID},
		"text":    {msg},
	}

	resp, err := http.PostForm(fmt.Sprintf("%s/sendMessage", t.APIURL), data)
	if err != nil {
		log.Printf("[ERROR] failed to send notificiation via telegram: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] failed to send notification via telegram: %v", resp.Status)
	}
}
