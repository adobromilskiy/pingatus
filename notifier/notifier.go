package notifier

import (
	"fmt"
	"log"
	"sync"

	"github.com/adobromilskiy/pingatus/config"
)

type Notifier interface {
	Send(msg string)
}

var (
	notifier      Notifier
	notifierError error
	notifierOnce  sync.Once
)

func Get(cfg *config.Config) (Notifier, error) {
	notifierOnce.Do(func() {
		switch cfg.Notifier.Type {
		case "telegram":
			if len(cfg.Notifier.TgToken) == 0 {
				notifierError = fmt.Errorf("telegram token is empty")
				return
			}
			if len(cfg.Notifier.TgChatID) == 0 {
				notifierError = fmt.Errorf("telegram chat id is empty")
				return
			}
			notifier = NewTelegram(cfg.Notifier.TgToken, cfg.Notifier.TgChatID)
		default:
			notifierError = fmt.Errorf("unknown notifier type: %s", cfg.Notifier.Type)
			return
		}

		log.Println("[INFO] notifier initialized")
	})

	return notifier, notifierError
}
