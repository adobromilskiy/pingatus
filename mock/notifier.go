package mock

import (
	"fmt"
)

type NotifierMock struct{}

func (n *NotifierMock) Send(msg string) {
	fmt.Printf("Send message: %s", msg)
}
