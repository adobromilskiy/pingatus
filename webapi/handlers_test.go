package webapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adobromilskiy/pingatus/mock"
)

func TestHandlerGetCurrentStatus(t *testing.T) {
	req, err := http.NewRequest("GET", "/status?name=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{
		Store: &mock.StoreMock{},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.handlerGetCurrentStatus)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandlerGetCurrentStatus_NoName(t *testing.T) {
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.handlerGetCurrentStatus)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := "name is required\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandlerGetCurrentStatus_Error(t *testing.T) {
	req, err := http.NewRequest("GET", "/status?name=test1", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{
		Store: &mock.StoreMock{},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.handlerGetCurrentStatus)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestHandlerGet24hrStats(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats?name=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{
		Store: &mock.StoreMock{},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.handlerGet24hrStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandlerGet24hrStats_NoName(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.handlerGet24hrStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := "name is required\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandlerGet24hrStats_Error(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats?name=test1", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{
		Store: &mock.StoreMock{},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.handlerGet24hrStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestHandlerGetEndpoints(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/endpoints", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{
		Store: &mock.StoreMock{},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.handlerGetNames)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `["test1","test2"]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
