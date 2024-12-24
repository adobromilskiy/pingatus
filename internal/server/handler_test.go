package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adobromilskiy/pingatus/internal/mocks"
)

func TestHandlerGetEndpoints(t *testing.T) {
	req, err := http.NewRequest("GET", "/endpoints", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{
		db: &mocks.StorageMock{},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.getEndpoints)

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

func TestGetEndpointStats(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats?name=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{
		db: &mocks.StorageMock{},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.getEndpointStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetEndpointStats_NoName(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.getEndpointStats)

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

func TestGetEndpointStats_Error(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats?name=test1", nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := &Server{
		lg: mocks.MockLogger(),
		db: &mocks.StorageMock{},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.getEndpointStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}
