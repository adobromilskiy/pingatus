package notifier

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewTelegram(t *testing.T) {
	logger := mockLogger()
	token := "test-token"
	chatID := "test-chat-id"

	tg := newTelegram(logger, token, chatID)

	if tg.token != token {
		t.Errorf("wrong token: got %v, want %v", tg.token, token)
	}

	if tg.chatID != chatID {
		t.Errorf("wrong chatID: got %v, want %v", tg.chatID, chatID)
	}

	expectedAPIURL := fmt.Sprintf("https://api.telegram.org/bot%s", token)
	if tg.apiURL != expectedAPIURL {
		t.Errorf("wrong apiURL: got %v, want %v", tg.apiURL, expectedAPIURL)
	}
}

func TestTelegramSend(t *testing.T) {
	tg := &telegram{
		lg:     mockLogger(),
		token:  "mock-token",
		chatID: "mock-chat-id",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/botmock-token/sendMessage" {
			t.Errorf("wrong URL path: got %v, want %v", r.URL.Path, "/botmock-token/sendMessage")
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}

		if chatID := r.Form.Get("chat_id"); chatID != tg.chatID {
			t.Errorf("wrong chat_id: got %v, want %v", chatID, tg.chatID)
		}

		if text := r.Form.Get("text"); text != "test message" {
			t.Errorf("wrong text: got %v, want %v", text, "test message")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tg.apiURL = server.URL + "/botmock-token"

	tg.Send(context.Background(), "test message")
}

func TestTelegramSend_Error(t *testing.T) {
	tg := &telegram{
		lg:     mockLogger(),
		token:  "mock-token",
		chatID: "mock-chat-id",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	tg.apiURL = server.URL + "/botmock-token"

	tg.Send(context.Background(), "test message")
}

func TestTelegramSend_ErrorPostRequest(t *testing.T) {
	tg := &telegram{
		lg:     mockLogger(),
		token:  "mock-token",
		chatID: "mock-chat-id",
	}

	tg.Send(context.Background(), "test message")
}
