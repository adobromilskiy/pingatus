package mocks

import (
	"errors"

	"github.com/adobromilskiy/pingatus/core"
)

type StorageMock struct{}

func (s *StorageMock) Save(data core.Endpoint) error {
	if data.Name != "test" {
		return errors.New("test: wrong name") //nolint:goerr113
	}

	return nil
}
