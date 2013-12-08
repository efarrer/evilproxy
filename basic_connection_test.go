package evil_proxy

import (
	"testing"
)

func TestConnectionBehaviorForBasicConnection(t *testing.T) {
	PerformConnectionTests(func() (Connection, Connection) {
		return NewBasicConnections(NewBasicPipe(), NewBasicPipe())
	}, t)
}
