package mocks

import "github.com/adobromilskiy/pingatus/core"

// Keep mocks referenced outside tests so deadcode doesn't flag them.
var (
	_                    = MockLogger
	_                    = (&NotifierMock{}).Send
	_ core.EndpointStore = (*StorageMock)(nil)
)
