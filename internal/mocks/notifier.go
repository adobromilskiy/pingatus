package mocks

import (
	"context"
	"fmt"
)

type NotifierMock struct{}

func (n *NotifierMock) Send(_ context.Context, msg string) {
	fmt.Printf("Send message: %s", msg)
}
