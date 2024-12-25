package mocks

import (
	"context"
	"errors"

	"github.com/adobromilskiy/pingatus/core"
)

type StorageMock struct{}

func (s *StorageMock) Save(_ context.Context, data core.Endpoint) error {
	if data.Name != "test" {
		return errors.New("test: wrong name") //nolint:goerr113
	}

	return nil
}

func (s *StorageMock) GetEndpoints(_ context.Context) ([]string, error) {
	return []string{"test1", "test2"}, nil
}

func (s *StorageMock) GetEndpointStats(_ context.Context, name string, date int64) ([]core.Endpoint, error) {
	if name != "test" {
		return nil, errors.New("test: wrong name") //nolint:goerr113
	}

	return []core.Endpoint{
		{
			Name:    name,
			Address: "http://localhost",
			Status:  true,
			Date:    date,
		}, {
			Name:    name,
			Address: "http://localhost",
			Status:  true,
			Date:    date,
		},
	}, nil
}
