package webapi

import (
	"testing"
	"time"

	"github.com/adobromilskiy/pingatus/storage"
)

func TestConvert(t *testing.T) {
	s := &Stats{}
	endpoints := []*storage.Endpoint{
		{
			Name:   "test",
			URL:    "http://test.com",
			Date:   time.Now().Unix(),
			Status: true,
		},
		{
			Name:   "test",
			URL:    "http://test.com",
			Date:   time.Now().Add(-1 * time.Hour).Unix(),
			Status: false,
		},
	}

	s.Convert(endpoints)

	if s.Name != "test" {
		t.Errorf("Expected %v, got %v", "test", s.Name)
	}

	if s.URL != "http://test.com" {
		t.Errorf("Expected %v, got %v", "http://test.com", s.URL)
	}

	if len(s.Hours) != 2 {
		t.Errorf("Expected %v, got %v", 2, len(s.Hours))
	}

	if len(s.Points) != 2 {
		t.Errorf("Expected %v, got %v", 2, len(s.Points))
	}
}

func TestConvert_empty(t *testing.T) {
	s := &Stats{}
	endpoints := []*storage.Endpoint{}

	s.Convert(endpoints)

	if s.Name != "" {
		t.Errorf("Expected %v, got %v", "", s.Name)
	}

	if s.URL != "" {
		t.Errorf("Expected %v, got %v", "", s.URL)
	}

	if len(s.Hours) != 0 {
		t.Errorf("Expected %v, got %v", 0, len(s.Hours))
	}

	if len(s.Points) != 0 {
		t.Errorf("Expected %v, got %v", 0, len(s.Points))
	}
}

func TestDuration_String(t *testing.T) {
	d := Duration(857)

	if d.String() != "14:17" {
		t.Errorf("Expected %v, got %v", "14:17", d.String())
	}

	d = Duration(0)

	if d.String() != "0:0" {
		t.Errorf("Expected %v, got %v", "0:0", d.String())
	}

	d = Duration(60)
	if d.String() != "1:0" {
		t.Errorf("Expected %v, got %v", "1:0", d.String())
	}

}
