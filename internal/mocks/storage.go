package mocks

import (
	"fmt"

	"github.com/adobromilskiy/pingatus/core"
)

type StorageMock struct{}

func (s *StorageMock) Save(data core.Endpoint) error {
	if data.Name != "test" {
		return fmt.Errorf("test: wrong name")
	}

	return nil
}
