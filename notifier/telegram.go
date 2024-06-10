package notifier

import "fmt"

type Telegram struct {
	Token  string
	ChatID string
}

func NewTelegram(token, chatID string) *Telegram {
	return &Telegram{token, chatID}
}

func (t *Telegram) Send(msg string) error {
	fmt.Println("Telegram: ", msg)
	return nil
}
