package mocks

import (
	"context"
	"errors"

	"github.com/adobromilskiy/pingatus/core"
)

type StorageMock struct{}

const testName = "test"

var errWrongName = errors.New("test: wrong name")

func (s *StorageMock) Save(_ context.Context, data core.Endpoint) error {
	if data.Name != testName {
		return errWrongName
	}

	return nil
}

func (s *StorageMock) GetEndpoints(_ context.Context) ([]string, error) {
	return []string{"test1", "test2"}, nil
}

func (s *StorageMock) GetEndpointStats(_ context.Context, name string, from int64, to int64) ([]core.Endpoint, error) {
	if name != testName {
		return nil, errWrongName
	}

	return []core.Endpoint{
		{
			Name:    name,
			Address: "http://localhost",
			Status:  true,
			Date:    from,
		}, {
			Name:    name,
			Address: "http://localhost",
			Status:  true,
			Date:    to,
		},
	}, nil
}
