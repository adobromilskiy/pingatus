package mock

import (
	"context"
	"fmt"

	"github.com/adobromilskiy/pingatus/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoreMock struct{}

func (s *StoreMock) GetLastEndpoint(_ context.Context, filter primitive.M) (*storage.Endpoint, error) {
	if filter["name"] != "test" {
		return nil, fmt.Errorf("test: wrong name")
	}
	enpoint := &storage.Endpoint{
		Name:   "test",
		Status: true,
		Date:   1234567890,
	}

	return enpoint, nil
}

func (s *StoreMock) GetEndpoints(_ context.Context, filter primitive.M) ([]*storage.Endpoint, error) {
	if filter["name"] != "test" {
		return nil, fmt.Errorf("test: wrong name")
	}

	endpoints := []*storage.Endpoint{
		{
			Name:   "test",
			Status: true,
			Date:   1234567890,
		},
	}

	return endpoints, nil
}

func (s *StoreMock) SaveEndpoint(_ context.Context, data *storage.Endpoint) error {
	if data.Name != "test" {
		return fmt.Errorf("test: wrong name")
	}
	return nil
}

func (s *StoreMock) GetNames(_ context.Context) ([]string, error) {
	return []string{"test1", "test2"}, nil
}

func (s *StoreMock) Close(_ context.Context) {}
