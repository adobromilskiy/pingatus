package notifier

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewTelegram(t *testing.T) {
	token := "test-token"
	chatID := "test-chat-id"
	tg := NewTelegram(token, chatID)

	if tg.Token != token {
		t.Errorf("wrong Token: got %v, want %v", tg.Token, token)
	}

	if tg.ChatID != chatID {
		t.Errorf("wrong ChatID: got %v, want %v", tg.ChatID, chatID)
	}

	expectedAPIURL := fmt.Sprintf("https://api.telegram.org/bot%s", token)
	if tg.APIURL != expectedAPIURL {
		t.Errorf("wrong APIURL: got %v, want %v", tg.APIURL, expectedAPIURL)
	}
}

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

func TestTelegramSendErrorPostRequest(t *testing.T) {
	tg := &Telegram{
		Token:  "mock-token",
		ChatID: "mock-chat-id",
	}

	tg.Send("test message")
}
