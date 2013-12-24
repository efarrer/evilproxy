package connection

import (
	"evilproxy/pipe"
	"testing"
)

func TestConnectionBehaviorForBasicConnection(t *testing.T) {
	PerformConnectionTests(func() (Connection, Connection) {
		return NewBasicConnections(pipe.NewBasicPipe(), pipe.NewBasicPipe())
	}, t)
}
