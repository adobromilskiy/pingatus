package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTelegramSend(t *testing.T) {
	tg := &Telegram{
		Token:  "mock-token",
		ChatID: "mock-chat-id",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/botmock-token/sendMessage" {
			t.Errorf("wrong URL path: got %v, want %v", r.URL.Path, "/botmock-token/sendMessage")
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}

		if chatID := r.Form.Get("chat_id"); chatID != tg.ChatID {
			t.Errorf("wrong chat_id: got %v, want %v", chatID, tg.ChatID)
		}

		if text := r.Form.Get("text"); text != "test message" {
			t.Errorf("wrong text: got %v, want %v", text, "test message")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// tg.apiURL = server.URL + "/botmock-token/sendMessage"
	tg.APIURL = server.URL + "/botmock-token"

	tg.Send("test message")
}

func TestTelegramSendError(t *testing.T) {
	tg := &Telegram{
		Token:  "mock-token",
		ChatID: "mock-chat-id",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	tg.APIURL = server.URL + "/botmock-token"

	tg.Send("test message")
}
